# Photo Grouping

Photo Grouping is a CLI tool that is able to generate photo group titles by location and time. It makes use of Google's 
[Geocoding Service](https://developers.google.com/maps/documentation/geocoding) in order to find the geological location
of where a photo was taken.

## What do I need to use this
In order to use Photo Grouping, a CSV file and API Key is required. 

The CSV format should be without headers and should specify the timestamp of the photo, along with the latitude and 
longitude, in that order. It should look something like this:
```
2022-04-01T18:52:59Z,40.627883,14.366858
2022-04-02T18:52:59Z,40.627883,14.366858
2022-04-03T18:52:59Z,40.627883,14.366858
```

The API Key is required is so the CLI is able to communicate with Google's Geocoding Service. In order to set up geocoding 
in your Google Cloud Project, the [following](https://developers.google.com/maps/documentation/geocoding/cloud-setup) 
documentation is available.

## Using the CLI
Once you have a CSV file and API Key, running the CLI can be done like so:
```
go run cmd/main.go --apiKey="<your_api_key>" --csvPath="path_to_your_csv.csv" 
```

After the command has been run, the output should look something like this:
```
A trip away to New York                      
New York in March                            
Visiting New York in March                   
A trip away to United States                 
United States in March                       
Visiting United States in March
```

## How does it work?
In order to determine titles for a group of photos, three factors are taken into consideration; 
* The location of the photo
* The first timestamp of a photo in a location
* The last timestamp of a photo in a location

Each row from the CSV is processed and the geological data is gathered. From there, each photo is categorised into groups 
based off of location (i.e London). Once grouped, timestamps can be used to figure out if the photos were taken on a 
weekend, during the week, a day trip, or a holiday.

### Assumptions
In order to give titles based on duration, some assumptions are made;
 1. day trip titles are only given to a group in which all photos were taken in the same day.
 2. week trip titles are only given to a group which spans between two and four days.
 3. weekend trip titles are given within the same period of a week trip but the dates must fall between Friday and Monday.
 4. holiday trip titles are only given to a group if the period is longer than four days.