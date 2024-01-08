const container = document.querySelector('.container');
const navBar = document.createElement('header');
const footer = document.createElement('footer');

navBar.innerHTML = `
<div class="navbar-container">
    <div class="navbar-left-side">
        <!-- <div class="navbar-pic">
            <img src="" alt="" height="60px" width="200px">
        </div> -->
        <div class="navbar-links">
            <a class="navbar-link-style" href="../pages/index.html">Home</a>
            <a class="navbar-link-style" href="../pages/about.html">About</a>
            <a class="navbar-link-style" href="../pages/donate.html">Donate</a>
        </div>
    </div>
    
        <div class="navbar-reg">
            <a class="navbar-link-style" href="../pages/login.html">Login</a>
            <a class="navbar-link-style" href="../pages/register.html">Sign Up</a>
        </div>
</div>
`

footer.innerHTML = `<footer>
<div class="footer-container">
    <div class="footer-left-side">
        <a href="" class="footer-link-style">Lorem, ipsum dolor.</a>
        <a href="" class="footer-link-style">Lorem, ipsum dolor.</a>
        <a href="" class="footer-link-style">Lorem, ipsum dolor.</a>
        <a href="" class="footer-link-style">Lorem, ipsum dolor.</a>
    </div>
    <div class="footer-left-side">
        <a href="" class="footer-link-style">Lorem, ipsum dolor.</a>
        <a href="" class="footer-link-style">Lorem, ipsum dolor.</a>
        <a href="" class="footer-link-style">Lorem, ipsum dolor.</a>
        <a href="" class="footer-link-style">Lorem, ipsum dolor.</a>
    </div>
    <div class="footer-left-side">
        <a href="" class="footer-link-style">Lorem, ipsum dolor.</a>
        <a href="" class="footer-link-style">Lorem, ipsum dolor.</a>
        <a href="" class="footer-link-style">Lorem, ipsum dolor.</a>
        <a href="" class="footer-link-style">Lorem, ipsum dolor.</a>
    </div>
</div>
</footer>
`

document.body.prepend(navBar);
container.after(footer);