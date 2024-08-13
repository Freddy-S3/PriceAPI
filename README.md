# PriceAPI

# Table of contents
1. [Introduction](#introduction)
2. [The problem](#theProblem)
3. [How to use](#howToUse)
4. [What I learned](#whatILearned)
5. [What I ran out of time for](#whatIRanOutOfTimeFor)
6. [Some comments](#someComments)



## Introduction <a name="introduction"></a>
Hi, my name is `Freddy Shaikh`,  
Thank you for taking the time to review me as an applicant!  
If you have any further questions, feel free to contact me at shaikhfh1@gmail.com


## The problem <a name="theProblem"></a>
I decided it was better not to link or describe the problem for privacy reasons.
I did however, leave the repo as public in case there were issues in the file submission process.


## How to use
To test functionality, simply run the main.go file from your IDE of choice (I used VSCode),  
and go to: http://localhost:5000/rates in your browser to access the `rates API`,  
and go to: http://localhost:5000/price in your browser to access the `price API`. (Will require a query)  

The `rates API` can take either a `GET` request which returns the priceDB.json file as a Http response in JSON format,
or a `PUT` request with the input being JSON data for new pricing info, which will save to priceDB.json.
The default seeded values will be loaded in upon server start so `PUT` data will be lost upon server restart.

The `price API` accepts queries such as `?start=2015-07-01T07:00:00-05:00&end=2015-07-01T12:22:00-05:00`,
and returns the expected price for the duration given the hourly rate within the database.  
Example Query: `http://localhost:5000/price?start=2015-07-01T07:00:00-05:00&end=2015-07-01T12:00:00-05:00`  
Expected Return in JSON: `{"price":8750}` 


## What I learned <a name="whatILearned"></a>
- While I have used Go before for some smaller projects and even a proof of concept, I have never used it to make and run a server which has been an interesting experience. I learned a lot about the `"net/http"` package. I learned about, downloaded, and used the Postman program for sending http `PUT` requests when manully testing to see if the http responses were correct. 

- I have done some basic JSON marshaling and unmarshaling before, but it was interesting to learn about encode and decode and how they require streams of data instead of the entire input at once. While I initially used marshal and unmarshal, I made the switch to encode and decode.

- After learning about "test tables" in Go I tried my best to implement them where possible. It's a great way to input multiple cases to see if you get the expected result. Admittedly, this requires more thinking of edge cases.

- At some point or another, I ended up reading parts of the documentation for every single package I imported. I had a window just for reading through packages. Became very close friends with the `"time"` package.


## What I ran out of time for <a name="whatIRanOutOfTimeFor"></a>
Unfortunately I ran out of time before I was able to implement:
- Unit tests for every function, integration tests, and end-to-end tests. Regretfully, I was unable to debug some of the HTTP tests I made, but I decided to leave them in so you can see what I was going for.
- Potentially refactor to make better use of `map[string]interface{}` in the decoding of the JSON, and potentially add concurrency to the database lookup in the `price API`. 
- A better formatted README


## Some comments <a name="someComments"></a>
- I intentionally used `net/http` instead of `gorilla/mux` for the handling of HTTP requests. I wanted to be sure to showcase that I can work with just the base packages for the sake of the assignment. However, if I were to continue building on this, or if this was not an assigment, I would use some framework like `gorilla/mux` to better handle server calls.

- I assumed that the rate used in the data base is hourly, so `price API` will return the total cost for the requested time band. If I am supposed to just return a flat rate instead of total for an hourly rate, I could just alter lines 101-108 in price.go. Initially, I actually implemented the flat rate method.

- For requests to the `price API`, I noticed that we can query up to the seconds, however the database for rate lookup only goes up to minutes. I just assumed theres a hard cutoff of zero seconds on the database rates, but I wanted to ask if I should accept queries up to the end of the minute? For example, if a rate it 0900-1000, do I accept a query ending on 10:00:59 or only up till 9:59:59 (as implemented).

- For requests to the `price API`, I understand I am to return `unavailable` at input spanning more than 1 day, but I was wondering if this meant more than 1 day respective to the timezone within the data base, or more than 1 day as a raw input (as implemented). For example, a request of `?start=2015-07-01T23:00:00-05:00&end=2015-07-02T00:01:00-05:00` of would return "unavailable" regardless of the timezone used in the rate within the database, even if it fit within the database time. 

- I assumed that the input provided during the `PUT` request to the `rates API` has no overlap in times so I didn't write any code to check for overlap from the rates. If I was required to create a check for the `PUT` request then I apologize.

- Although Python is my stronger language, I chose Go knowing that it was faster, and better at implementing APIs among other reasons (concurrency, learning experience, for fun).

- When writing the data structure to use in the price lookup, I considered writing a map that used 3 keys: "Day, StartTime, EndTime" and resulted in 1 value: "Price" This would create a more linear lookup time, however, this would also require significantly more space in memory. To reduce the space in that idea, I considered taking each DB entry and registering it to a weekly calendar of sorts where each day would have a list of timeband but this would have required a lookup of the registered DB time bands for a single day regardless. Ultimately I decided against it because it would alter the structure to the point where the JSON DB wouldn't make as much sense. I was considering revisiting this idea in refactoring, but ultimately I ran out of time.
