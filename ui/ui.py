import base64
import json
import os
import io

import requests
import streamlit as st
import chess
import chess.pgn
import chess.svg


#---------------------#
# Streamlit setup     #
#---------------------#

st.set_page_config(layout="wide")

st.title("LLaMA 3 Chess")
col1, col2 = st.columns(2)


#---------------------#
# Board setup         #
#---------------------#

if "board" not in st.session_state:
    #pgn = open('games/jimmy1400.txt')
    #game = chess.pgn.read_game(pgn)
    game = chess.pgn.Game.from_board(chess.Board())
    board = game.board()
    for move in game.mainline_moves():
        board.push(move)
    st.session_state["board"] = board

def update_game(updated_game):
    game = chess.pgn.read_game(io.StringIO(updated_game))
    board = game.board()
    for move in game.mainline_moves():
        board.push(move)
    st.session_state["board"] = board

def render_svg(placeholder, svg):
    """Renders the given svg string."""
    b64 = base64.b64encode(svg.encode('utf-8')).decode("utf-8")
    html = r'<img src="data:image/svg+xml;base64,%s" width="500px" align="center" style="padding-bottom: 20px; display: block; border: none;"/>' % b64
    placeholder.write(html, unsafe_allow_html=True)


#---------------------#
# LLM chess API setup #
#---------------------#

url = os.getenv("CHESS_API_URL")

def parse_game_pgn(pgn_raw):
    updated_game = pgn_raw
    updated_game = updated_game.split('\n\n')[-1].strip()
    return updated_game

def generate_move():
    payload = json.dumps({
        "game": parse_game_pgn(str(chess.pgn.Game().from_board(st.session_state["board"])))
    })
    headers = {
        'Content-Type': 'application/json'
    }
    response = requests.request("POST", url + "/move", headers=headers, data=payload).json()
    updated_game = parse_game_pgn(response['game_updated'])
    return updated_game

def parse_move(move_text):
    payload = json.dumps({
        "game": parse_game_pgn(str(chess.pgn.Game().from_board(st.session_state["board"]))),
        "move": move_text
    })
    headers = {
        'Content-Type': 'application/json'
    }
    response = requests.request("POST", url + "/parse", headers=headers, data=payload).json()
    updated_game = parse_game_pgn(response['game_updated'])
    return updated_game

def get_help():
    payload = json.dumps({
        "game": parse_game_pgn(str(chess.pgn.Game().from_board(st.session_state["board"]))),
    })
    headers = {
        'Content-Type': 'application/json'
    }
    response = requests.request("POST", url + "/help", headers=headers, data=payload).json()
    advice = response['message']
    return advice


#---------------------#
# Play the game       #
#---------------------#

with col1:
    placeholder = st.empty()
    render_svg(placeholder, chess.svg.board(st.session_state["board"]))

    # Check if there is a game outcome.
    if st.session_state["board"].outcome() is not None:
        st.balloons()
        if st.session_state["board"].outcome().winner == chess.WHITE:
            st.markdown("**Game status:** You won! You are smarter than a LLaMA.")
        elif st.session_state["board"].outcome().winner == chess.BLACK:
            st.markdown("**Game status:** LLaMA 3 won!")
        else:
            st.markdown("**Game status:** Draw. Are you a LLaMA?")
    else:
        st.markdown("**Game status:** In progress. Use the boxes to the right to play.")

with col2:
    st.markdown("### Make a move")
    move_text = st.text_input(
        label="Type somthing like \"Pawn to c4\", \"bishop takes queen at h8\", \"Pawn on c-file takes pawn at b5\" or \"Knight to d5\".", 
        placeholder="Pawn to e4",
        key="move"
    )
    if submit := st.button("Move"):
        if move_text:

            with st.spinner("Parsing your move..."):

                understood = True
                try:

                    # Parse the move.
                    game_updated = parse_move(move_text)

                    # Update the game.
                    update_game(game_updated)
                    render_svg(placeholder, chess.svg.board(st.session_state["board"]))

                except:

                    understood = False
                    st.error("Either the move is invalid or LLaMA 3 didn't understand it. Please try again or be more specific.")

            if understood:
                with st.spinner("Generating LLaMA3's next move..."):

                    try:

                        # Let the LLM generate a move.
                        game_updated = generate_move()

                        # Update the game.
                        update_game(game_updated)
                        render_svg(placeholder, chess.svg.board(st.session_state["board"]))

                    except:

                        st.error("LLaMA 3 got stumped and gave up. You are smarter than a language model. Refresh to play again.")

                st.success("Your move was made and LLaMA 3 responded! Move again.")

        else:
            st.warning("Please enter a move.")

    st.divider()
    st.markdown("### Need help?")
    if submit := st.button("Get help"):
        with st.spinner("Analyzing the game with LLaMA 3 power..."):
            advice = get_help()
            st.markdown(advice)
