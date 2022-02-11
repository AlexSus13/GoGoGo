# GoGoGo
1) To register a user, send a request in the form:

curl -v -X POST -d '{"user_name": "yourname", "user_password": "yourpassword"}' http:/37.139.43.30:8080/reg

2) To authenticate a user, send a request in the form:

curl -v -X POST -d '{"user_name": "yourname", "user_password": "yourpassword"}' http:/37.139.43.30:8080/auth

3) For information about all files, send a request in the form:

curl -v -X GET -H "Authorization:Bearer YourToken" http:/37.139.43.30:8080/api/listfiles

4) To get information about files by part of the name, send a request in the form:

curl -v -X GET -H "Authorization:Bearer YourToken" http:/37.139.43.30:8080/api/listfiles?filename=PartOfTheName

5) To get the contents of the file, send a request in the form:

curl -v -X GET -H "Authorization:Bearer YourToken" http:/37.139.43.30:8080/api/getfile?filename=FileName.txt

6) To add a file, send a request in the form:

curl -v -X POST -d @/PathToTheAddedFile/FileName.txt -H "Authorization:Bearer YourToken" http:/37.139.43.30:8080/api/postfile?filename=FileName.txt

7) Example .env file you can look in SampleEnvFile.txt
