# Thugger Lyrics Service

Golang service that serves as an HTTP API for Thugger Lyrics.

## Install

```sh
git clone git@github.com:jseashell/thugger-lyrics-service
cd thugger-lyrics-service
make
```

## Endpoints

> DynamoDB can be seeded with [https://github.com/jseashell/aws-genius-lyrics](https://github.com/jseashell/aws-genius-lyrics)

### /random

Fetches a random Young Thug song from DynamoDB

```sh
curl https://<api-id>.execute-api.us-east-1.amazonaws.com/random
{"artist_names":"","full_title":"","header_image_thumbnail_url":"","header_image_url":"","id":0,"uuid":"","path":"/Young-thug-wake-up-high-lyrics","release_date_for_display":"","song_art_image_thumbnail_url":"","song_art_image_url":"","title":"","url":"https://genius.com/Young-thug-wake-up-high-lyrics","lyrics":["Every day I wake up high, go to sleep, wake up, do it again","I done made the spreadsheet for all my homies like they kin","I done dodged the [?], I done dodged the questions","I dodged the police 'cause they don't know nothing 'bout this new Benz","I called to the street 'cause I woke up to one of my friends shot","Call the police, yeah","Pouring codeine, yeah","Lil' shawty clean, yeah","Mismatch my diamonds, blinding","I don't want it 'less you tell me I'm charming","I don't want it 'less I see you trying","Ain't giving hеr up, nobody seeing you crying","Wrap you up, they ain't seeing your face","It's a Berеtta on my lil' bitty waist","I get that cream and that cheddar every base","She waking up and taking shots right to the face","With GPS, the car come right to where you stay","Your silky body looking just like suede","You my angel, I don't want nobody else, baby","Peep my rasta, I been all exotic","Every day I wake up high, go to sleep, wake up, do it again","I done made the spreadsheet for all my homies like they kin","I done dodged the [?], I done dodged the questions","I dodged the police 'cause they don't know nothing 'bout this new Benz","I called to the street 'cause I woke up to one of my friends shot","Call the police, yeah","Pouring codeine, yeah","Lil' shawty clean"]}
```

## Deploy

1. [Configure the AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html) on your local workstation.
1. Run the deploy script

    ```sh
    make deploy
    ```

## Project Structure

```text
.
├── random                  # package for /random endpoint
│   └── main.go
├── .gitignore
├── go.mod                  # module dependencies
├── go.sum                  # dependency checksums
├── LICENSE
├── Makefile                # build script
├── README.md
└── serverless.yml          # serverless framework configuration
```

## 3rd party libraries

- [aws-lambda-go](https://github.com/aws/aws-lambda-go) - AWS Lambda SDK for the Go programming language.
- [google/uuid](https://github.com/google/uuid) -  RFC-4122 compliant UUID module by Google.

## Disclaimer

Repository contributors are not responsible for costs incurred by AWS services.

## License

This software is distributed under the terms of the [MIT License](/LICENSE)