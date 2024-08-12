# Housemate Expense Manager

I was working on an expense management coding challenge and found it fascinating how Splitwise simplifies transactions. Inspired by this, I implemented a similar solution in Go, focusing on minimizing the number of transactions among housemates.

## Overview

The Housemate Expense Manager allows up to three members in a house to track shared expenses, settle dues, and manage payments efficiently.

### Key Commands

- **MOVE_IN `<name>`**: Adds a member to the house. Returns `SUCCESS` if successful or `HOUSEFUL` if the house is full.

- **SPEND `<amount>` `<spent-by>` `<spent-for...>`**: Tracks expenses shared among specified members. Returns `SUCCESS` or `MEMBER_NOT_FOUND` if any member is missing.

- **DUES `<member>`**: Displays all outstanding dues for a member, sorted by amount and name.

- **CLEAR_DUE `<payer>` `<payee>` `<amount>`**: Allows a member to clear their dues. Returns the remaining balance or `INCORRECT_PAYMENT` if the payment exceeds the owed amount.

- **MOVE_OUT `<name>`**: Allows a member to move out if all dues are settled. Returns `SUCCESS`, `FAILURE` if dues remain, or `MEMBER_NOT_FOUND` if the member doesn't exist.

### Example Usage

```plaintext
MOVE_IN ALICE
MOVE_IN BOB
MOVE_IN CHARLIE
MOVE_IN DAVID
SPEND 3000 ALICE BOB CHARLIE
SPEND 300 BOB CHARLIE
SPEND 300 BOB DAVID
DUES CHARLIE
DUES BOB
CLEAR_DUE CHARLIE ALICE 500
CLEAR_DUE CHARLIE ALICE 2500
MOVE_OUT ALICE
MOVE_OUT BOB
MOVE_OUT CHARLIE
CLEAR_DUE CHARLIE ALICE 650
MOVE_OUT CHARLIE

OUTPUT:
SUCCESS
SUCCESS
SUCCESS
HOUSEFUL
SUCCESS
SUCCESS
MEMBER_NOT_FOUND
ALICE 1150
BOB 0
ALICE 850
CHARLIE 0
650
INCORRECT_PAYMENT
FAILURE
FAILURE
FAILURE
0
SUCCESS

