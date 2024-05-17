# CARBON - Logger Search Engine

CARBON is a lightweight tool designed to efficiently search through **OUR** log files using Golang. It utilizes an inverted index and n-grams to provide fast and accurate results, making it ideal for parsing and analyzing large log datasets. Our logs are comprised of words describing the action as opposed to the standard template `MAY 17 00:57:06.11 GET 200 localhost:8080 200 /product/view/product-name message`.

NOTE: It was designed to solve a specific problem at work but I am currently working on converting it to a CLI app for general purpose use, optimizing it to search faster through the template logs instead of being faster with text based search.

## Features

- **Efficient Search Algorithm:** Implements an inverted index and n-grams for quick log lookup.
- **Lightweight and Scalable:** Designed to handle large log files.
- **Easy to Use:** Simple interface for searching and analyzing logs. Just run it and make a GET to localhost:8080/search?q=

## Models

- **The Concurrency Model:** For nGram sizes bigger than 5, please use this model as it will be significantly faster than the other models.
- **The Lightweight Model:** This model doesn't use nGrams, it is used purely used for 1:1 word match, it doesn't autocomplete nor does it find substrings. It's the fastest and the lightest out of the 2 no matter the nGram sized used for the other one.

## Purpose

Originally developed to address a specific logging challenge I had at work, CARBON quickly became an invaluable tool for me and my team. Its effectiveness in improving log analysis workflows led to its adoption as a standard utility within my team.

## Usage

To use CARBON, simply clone the repository, build your prefered model and run it.

## Technologies Used

- Golang
