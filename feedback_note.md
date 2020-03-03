# Feedback

## Language
Score: 3/5 ==> Meets Expectations

1. Parsing environment variables and while ensuring required fields are available can be achieved with lesser boiler 
plate code by leveraging already existing  env vars parsers e.g (https://github.com/caarlos0/env)
https://github.com/wechuli/rss_feeds_fetcher/blob/master/main.go#L14

2. Server port should be configurable through environment variables. 
https://github.com/wechuli/rss_feeds_fetcher/blob/master/main.go#L19

3. Initializing the App should be encapsulate within the app package. Usually go standard practice would be to have an 
init function like app.New(<dependencies>) which returns a properly initialized App and then there wouldn't be need for 
https://github.com/wechuli/rss_feeds_fetcher/blob/master/main.go#L17. Or since App is part of package main
(https://github.com/wechuli/rss_feeds_fetcher/blob/master/app.go#L29) it shouldn't be a method of the App struct but a
standalone Initialize function which returns a properly initialized App.

4. Defering the Body#Close should be done immediately, and not before handling the error. Also, that function returns an
error which is ignored. Ideally we should report it atleast logging as a warning 
(https://github.com/wechuli/rss_feeds_fetcher/blob/master/app.go#L81)

5. Logic for Fetching the news feed and populating should have been made concurrent (ease to implement concurrency 
is one of the best features of Go in my opinnion :)). You can see https://blog.golang.org/pipelines for examples.
https://github.com/wechuli/rss_feeds_fetcher/blob/master/app.go#L92

6. DB creation should be separated from application logic. It's not a good idea to have the creation script ran everytime 
as this is not efficient and also not scalable if we have multiple Tables. See https://github.com/golang-migrate/migrate. 
The app should expect that the DB is already properly initialized.
https://github.com/wechuli/rss_feeds_fetcher/blob/master/models/model.go#L47

7. Uses Go modules

## Architecture and Design
Score: 3/5 ==> Meets Expectations

1. Any specific reason why /search endpoint is a POST endpoint. This ideally should be a GET  endpoint.

2. App could be made more robust by:
a. Checking that the required environment variables are available at start up 
(or using default values)

3. DB Port should also be configurable through environment variables
(https://github.com/wechuli/rss_feeds_fetcher/blob/master/main.go#L9)

4. Application logging  can definitely be improved by atleast using a strutured logger which supports different 
log levels (e.g https://godoc.org/go.uber.org/zap)

5. I do not believe we need use a relational database for the following reasons:
- Our datastore model will have only one Entity type. No relationships between several Entities
- If we used a Non-Relational DB engine like MongoDB (No-SQL) we would not need boiler plate code for migrations 
and the likes.

6. We should have taken advantage of Indexing to speed up DB search.

## Testing and Security
Score: 3/5 ==> Meets Expectations

1. Please note that I understand that these feedback might not necessarily be feasible given the implementation timelines
but still it probably should have been documented somewhere as TODOs (maybe in the Project README so we are certain you had them in mind)

2. We could implement RateLimiting to prevent potential DDOS attacks.

# Summary
## Overall Average Rating = 3/5
