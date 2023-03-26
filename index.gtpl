<!DOCTYPE html>
<html>
<head>
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<style>
* {
  box-sizing: border-box;
}

.menu {
  float: left;
  width: 20%;
  text-align: center;
}

.menu a {
  background-color: #e5e5e5;
  padding: 8px;
  margin-top: 7px;
  display: block;
  width: 100%;
  color: black;
}

.main {
  float: left;
  width: 60%;
  padding: 0 20px;
}

.right {
  background-color: #e5e5e5;
  float: left;
  width: 20%;
  padding: 15px;
  margin-top: 7px;
  text-align: center;
}

@media only screen and (max-width: 620px) {
  /* For mobile phones: */
  .menu, .main, .right {
    width: 100%;
  }
}
</style>
</head>
<body style="font-family:Verdana;color:#410000;">

<div style="background-color:#e5e5e5;padding:15px;text-align:center;">
  <h2>Welcome to my homework GUI</h2>
  <h2>Yeah,it is wierd to say that But I just said it.</h2>
</div>

<div style="overflow:auto">
  <div class="menu">
    <a href="/login">Login</a>
    <a href="/register">Sign in</a>
    <a href="/Encription">Homework</a>
    <a href="/Extra">Extra</a>
  </div>

  <div class="main">
    <h3>While considering the most appropriate programming language, we had some back and forth before finally deciding to use Go-lang. However, at one point, we asked ourselves, "Why not make it a server?"</h3>
    <p>Before using Go-lang, we had tried writing our code in various programming languages such as C++, C#, Rust, and Python.</p>
</div>

  <div class="right">
    <h2>About</h2>
    <p>Fahrettin Baştürk</p>
    <p>Turgut Baştürk</p>
  </div>
</div>

<div style="background-color:#e5e5e5;text-align:center;padding:10px;margin-top:7px;">© copyright Fahrettin Baştürk</div>

</body>
</html>

