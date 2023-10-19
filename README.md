# To Run this Project

- You need to have `Go v1.17` installed for the back end
- You'll also need `Angular CLI v12.1.2` and `Node v14.15.0` installed for the front end
- Clone this project
- Make sure you have a Local MongoDB database set up on port 27017
- To execute the back end of the app, run `go run app/main.go`
- On Postman execute the following cURL: `curl --location --request POST 'http://localhost:3000/v1/preload'` to populate the DB with the required Users and Chatrooms data (you should see 3 users under Users collection, and 2 chatrooms under Chatrooms collection)
- Git clone the front project from [this URL](https://github.com/ratolibre1/chat-ang)
- Run `ng serve` on the root of the front project

# To Review this Project

- Load `http://localhost:4200` on two browsers
- You should see the Login section
- Use credentials `email: artichoke@jobsity.com, password: 123456` on one of them and `email: lilbasil@jobsity.com, password: 123456` on the other (wrong email and/or password will result in failure to login)
- The Chatrooms section should appear
- Enter the same chatroom with both users. When you type a message with one of them it will show on both chats (unless it's a command)
- Chat autoscrolls when a new message is received (most of the time, it doesn't behave too well with the stock responses)
- You can change chatrooms
- If you send a valid stock command (like `/stock=aapl.us`) the message won't be stored but a response will be generated and stored

# Some considerations

- I tried to complete all the mandatory points. I got stuck implementing Unit Tests as it's something I've never done before in Golang (I can write UTs in Angular just fine)
- I chose heagonal architecture as I believe it is very orderly and allows devs to separate the layers of the app
- I tried to work as if I was part of a team doing Scrum. Instead of Jira I kept a [Trello Board](https://trello.com/invite/b/sgfmFApA/ATTI4dd11900533d406bf88694ea4daa3309FCC10A6A/chat-hex) with my list of tasks.
- I implemented JWT for Authentication as they are simple enough to implement and secure enough to block access to protected URLs
- The RabbitMQ code is a bit sloppy as I was focusing more in getting it done than in doing it nice. Apologies for that, it still gets the job done
- I wanted to implement websockets for the chat to work in real time. I didn't have the time to so instead the front just pings the messages endpoint once per second as an MVP
- The front is a bit ugly but I expect it to be understandable enoguh. I promise I can write more beautiful Angular code if given the time
  
- You might wonder the meaning behind the two test emails. We have two wonderful calico cats called Alcachofa ("Artichoke" in Spanish) and Albahaquita ("Little Basil" in Spanish). You can see them both [here](https://photos.app.goo.gl/aNXfAAouLwFQPUag7) (Chofi is the black headed one, Baqui is the white headed one)
- Thank you for the opportunity! I hope this project showcases my skills and what I can bring to the team in terms of coding
