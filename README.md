# Todo application
This is a simple todo application created in golang. The aim is to have a quick project in golang and make it live.

## Project Goals
1. Create a simple todo app which you can run from your terminal for yourself
2. Deploy it in a server hosted in cloud and make it accessible using an endpoint. 
3. Create a simple web page layer over this in simple javascript
4. Add user authentication part and seperate the design to behave per user.

## Tips
* Take notes while you are building about golang and any concepts in general. Make it as part of the project.
* Finish this is a high paced fixed time rather than it being an ongoing project. (probably over a single day over a weekend)

## Sources
https://betterprogramming.pub/build-a-simple-todolist-app-in-golang-82297ec25c7d
https://keiran.scot/2018/03/02/building-a-todo-api-with-golang-and-kubernetes-part-1-introduction/



## Commands
docker run -d -p 3306:3306 --name mysql -e MYSQL_ROOT_PASSWORD=root --platform linux/x86_64 mysql

## CRUD 
curl -i localhost:8000/healthz  
curl -i -X PUT -d 'description=Feed the cat 'localhost:8000/todo'  
curl -i -X POST -d 'completed=true 'localhost:8000/todo/1'  
curl -i -X DELETE 'localhost:8000/todo/2'  
curl -i -X GET 'localhost:8000/todo-completed'  
curl -i -X GET 'localhost:8000/todo-incomplete'  