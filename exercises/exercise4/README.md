# Go GenAI App Build - Exercise 3

The goal of this exercise is to generate helpful advice (using RAG) for the user. Add your code to the `TODO` secion in `llm.go` to accomplish this task. You can optionally setup your own Vector DB and load in data from a chess strategy book.

## Setup the Vector DB

In order to get the user helpful advice, we will integrate an external data source of chess information. Specifically, we will utilize Retrieval Augmented Generation (RAG) to search a chess strategy book and augment prompts to the LLM with information from this book. This will require us to parse and load the book into a database.

The instructor has already created a database and loaded the data into the database with the scripts in [../../db](../../db). Specifically, these scripts were run in the following order:

1. `parse` - Parse the book data
2. `embed` - Embed the chunks
3. `load` - Load the chunks into the database
4. `query` - Test a query 

If you want to replicate the database setup, you will need a Postgres database with the `vector` extension. An easy way to set this up is to create a free account on [Neon](https://neon.tech/). Create a new database, and then execute the following to setup the database:

```
CREATE EXTENSION vector;

CREATE TABLE items (
  id integer PRIMARY KEY,
  chunk varchar,
  metadata varchar,
  embedding VECTOR(512)
);
```

Once this is setup, you can follow the steps outlined above to parse and load in the book data.

## Complete the exercise

Execute the following steps:

1. Add your functionality to the `TODO` in `llm.go`
2. Build and run the updated API
3. Test the API health check with Postman (or similar) or via cURL:

```
curl --location 'http://localhost:8080/help' \
--header 'Content-Type: application/json' \
--data '{
    "game": "1. Nf3 d5 2. c4 Nc6 3. cxd5 Nf6 4. d4 Bg4 5. Bg5 e6 6. Bxf6 Qxd5 7. e4 Bb4+ 8. Ke2 Qd7 9. Nc3 O-O 10. a3 Be7 11. e5 Qd8 12. Ke1 *"
}'
```

Which should return something like:

```json
{
    "message": "In the current position, White has just played 12. Ke1. Black has a solid position with well-developed pieces and a castled king. White's strategy should be to maintain control of the center and look for opportunities to expand, while Black aims to create counterplay and challenge White's central pawns.\n\nFor the next move, consider:\n\n1. 12... Qd6: This move supports the central pawn on d5 and prepares for possible counterplay. It also puts pressure on White's pieces, particularly the knight on c3.\n\n2. 12... Nf6: This move aims to challenge White's central pawn on e4 and prepare for the development of the dark-squared bishop. It also puts pressure on White's knight on f3.\n\n3. 12... O-O: This move consolidates Black's position and prepares for further counterplay. It is a safe move that allows Black to focus on developing their pieces and looking for opportunities to create counterplay.\n\nIn general, Black should continue to develop their pieces and look for opportunities to create counterplay against White's central pawns. Avoiding unnecessary pawn moves and focusing on piece development will be crucial in the upcoming middle game.",
    "reference_info": "\r\nThe idea underlying this pawn sacrifice is to open the K file for\r\nthe Rook. It will be seen that, with correct play, Black manages\r\nto castle just in time, and White, though winning back his pawn,\r\nhas no advantage in position. The opening is seldom played by\r\nmodern masters.\r\n\r\nInstead of the move in the text, White can hardly defend the KP\r\nwith Kt-B3, as Black simply captures the pawn and recovers his\r\npiece by P-Q4, with a satisfactory position.  It is even better\r\nfor Black if White plays 6. BxPch in reply to 5. ... KtxP. The\r\ncapture of White’s KP is far more important than that of the\r\nBlack KBP, particularly as the White Bishop, which could be\r\ndangerous on the diagonal QR2-KKt8, is exchanged, e.g. 6. ...\r\nKxB; 7. KtxKt, P-Q4; 8. Kt-Kt5ch, K-Kt1! Black continues P-KR3,\r\nK-R2, R-B1 and has open lines for Rooks and Bishops.\r\n\r\n          4. ...             KtxP\r\n\r\nBlack can, of course, develop his B-B4. Then he must either\r\nsubmit to the Max Lange attack (5. P-Q4, PxP) or play BxP, giving\r\nup the useful B, in which case he loses the pawn gained after 6.\r\nKtxB, KtxKt; 7. P-KB4, P-Q3; 8. PxP, PxP; 9. B-KKt5, and\r\neventually Q-B3.\r\n\r\n          5. P-Q4\r\n\r\nR-K1 at once would lead to nothing.\r\n\r\n          5. ...                 PxP\r\n          6. R-K1                P-Q4\r\n          7. BxP!                QxB\r\n          8. Kt-B3\r\n\r\n        ---------------------------------------\r\n     8 | #R |    | #B |    | #K | #B |    | #R |\r\n       |---------------------------------------|\r\n     7 | #P | #P | #P |    |    | #P | #P | #P |\r\n       |---------------------------------------|\r\n     6 |    |    | #Kt|    |    |    |    |    |\r\n       |---------------------------------------|\r\n     5 |    |   "
}
```