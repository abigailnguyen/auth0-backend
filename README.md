## Project Idea
Learn how to build and secure a Go API with JSON Web Tokens (JWTs) and consume it with a modern ReactUI. Users will authenticate on the React side with Auth0 and then make a request to the go API by sending their access token along with the request.  

## Setup 
1. Setup an Auth0 account, by navigating to [Auth0 Sign up](https://auth0.com/signup)
2. Naviagte to management dashboard and **Create Application > Single Page Web Applications**
3. Click on **Settings** and fill in with `http://localhost:3000`      
- **Allowed Callback URLs** 
- **Allowed Logout URLs**
- **Allowed Web Origins** 

    and then select **Save Changes**
4. Initialize `go mod` - anything without http://, ideally it should be a repository URL in case the code is shared amongst different projects
    > go mod init github.com/abigailnguyen/go-website

5. For simplicity, all the application will be in the file `main.go`
6. To import the packages for used in `main.go`
    > go get -u "github.com/gorrila/mux"

   or for all dependencies
    > go get

7. To run the application
    > go run main.go

## Defining our API

1. Our API is going to consist of 3 routes:
    - status - will handle our call to make sure that our API is up and running
    - products - will retrieve a list of products that the user can leave feedback on
    - products/{slug}/feedback - will capture user feedback on products

2. In addition, we implement a handler function called `NotImplemented`, which will be the default handler for routes with no custom functionality yet.

3. Talking about handlers/middleware: **In Go, middleware is refereed to as handlers**. It is abstracted code that runs before the intended code is executed. For example, you may have a loggin middleware that logs information about each request. You wouldn't want to implement the logging code for each route individually, so you would write a middleware function that gets inserted before the main function of the route is called that would handle the logging. **We will use custom handlers to secure our API.**

## Front end - React
**An API is only as good as the frontend that consumes it**

1. Create a default React App 
    ```bash 
    npx create-react-app static
    cd static
    npm start
    ```
2. Install dependencies needed on the React side.
    ```bash
    npm install react-router-dom @auth0/auth0-spa-js bootstrap react-icons
    ```
    > - `@auth0/auth0-spa-js` - Auth0's JavaScript SDK for SPAs
    > - `react-router-dom` - React's router package
    > - `bootstrap` - For quick styling (optional)
    > - `react-icons` - Add icons to the app (optional)

3. Setting up components
Your application will allow the user to view and leave feedback on VR products. Users must be logged in before they can leave feedback.

You will ned to build 3 components:
- `App` - to launch the application
- `Home` - will be displayed for non-logged in users
- `LoggedIn` - will display the products available for review for authenticated users

First, create a new folder inside `src` called `components` and create 2 files 
```bash
cd static/src/components
touch Home.js LoggedIn.js
```

## Setting up Auth0
Because some of the components will be depending on the authentication state, go ahead and set up the Auth0 React wrapper. 
```bash
cd .. 
touch react-auth0-spa.js
```

1. **src/react-auth0-spa.js**: a wrapper that creates functions that easily integrates with the rest of the React components
to allow the user to login and logout and display user information.

2. **src/history.js**: this file helps redirect after the user signs in
 
3. **src/index.js**: Integrate Auth0 SDK into React application. This wraps the *\<App>* component inside *\<Auth0Provider>* component

4. You now need to fill in the *domain* and *client_id* values.
    **src/auth_config.json** : drop the values that you created in Auth0 for your application here


## App Component
The **\<App>** component is where your application starts.

## Home Component
The **\<Home>** component displays when the user is not yet logged in.

## LoggedIn Component
The **\<LoggedIn>** component will display after the user has a valid access token, and successfully logged in.
This could be split into multiple components but for simplicity we will keep them all in one.

## Authorization with Golang
Adding authorization will allow you protect your API. Since you are dealing with projects in development,
you do not want any data to be publicly available. You have accomplished the first step of logging in the React appplication.
The next step is to pull data from the Go application, but only if the user has a valid access token.

1. In **main.go**, you still have two handlers for POST and GET endpoints to get and update the products.
Since you are updating this application to use Auth0 for authorization and authentication, a new handler **jwtHandler**, 
has been added. This will wrap around the endpoints you want to protect. You also need to enable CORS to allow Cross-origins
requests so that the React client can consume the API

## Create the middleware
Next create the middleware that will be applied to your endpoints. This middleware will check if an access token exists and is valid.
If it passes the checks, the request will proceed. If not, a 401 Authorization error is returned.
The structs used for this middleware is saved inside file *response.go*

1. The **JSONWebKeys** struct holds fields related to the JSON Web Key Set for this API. These keys contain the public keys,
which will be used to verify JWTs.

2. You have the new ***auth0/go-jwt-middleware*** middleware function that will validate tokens coming from Auth0. There are some
values here that you'll need to fill in, so take note of those, and you'll come back to them once setup is finished.

3. Next you need to create the function to grab the JSON Web Key Set and return the certificate with the public key,
called ***getPemCert***

4. Now that the middleware is set up, you need to apply it to the endpoints you want to protect.

## Apply the middleware
Any private endpoints that you want to protect in the future should also use **jwtMiddleware**
    
    r.Handle("/products", jwtMiddleware.Handler(ProductsHandler)).Methods("GET")

## Add the imports
- **auth0/go-jwt-middleware** - Auth0 package that fetches your Auth0 public key and checks for JWTs on HTTP requests
- **rs/cors** - CORS is a net/http handler implementing CORS specification of Go

## Setting up the Auth0 API
Now you need to register your Go API with Auth0. 
1. Navigate to the APIs section, and create a new API by clicking the **Create API** button
2. Give your API a Name and an Identifier. These can be anything you'd like, but for the Identifier, it's best to use a URL format
for naming convention purposes. This doesn't have to be a publicly available URL, and Auth0 will never call it. You can leave the 
Signing Algorithm as is (RS256). Once you have that fileed out, click Create.
3. Grab the value for *options.Authority* and replace the variable *aud* value with it.
4. Copy the value of *options.Audience* and replace "https://YOUR_DOMAIN" with it for *iss* variable and *getPermCert* function.

## Testing with cURL

    curl --request GET \
      --url http://localhost:8080/products \
      --header 'authorization: Bearer yourtokenhere'
    
You can grab the token under the API -> Test tab -> Response

## Connecting React and Go
So you now have users able to sign in on the React side and API authorization implemented on the backend.
The final step is to connect the two. The goal here is that a user will sign into the frontend, and then you'll 
send their token to the backend with the request for product data. If their access token is valid, then you'll 
return the data.

To accomplish this, you only need to update the React side.
1. You're adding an audience value in **src/index.js**, so you have to update **auth_config.json** file with that value.
The audience value can be found by copy the Identifier value under Settings, and paste it there.
## Notes:
1. If you find that some of the imports disappear on save, make sure that your code editor isn't set to remove unused imports 
on save. In VSCode, open up your "User settings" and search for *gofmt*. If it's set to *goreturn* then change that to *gofmt*

2. You can run **go get** or **go build main.go** here, which will also grab the dependencies and compile

3. Re run your application with **go run .**

