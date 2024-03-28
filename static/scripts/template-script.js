"use strict"

const burgerButton = document.querySelector('.burger-btn');
const navBarContent = document.querySelector('.navbar-content');
const dropdownMenu = document?.querySelector('.dropdown');
const dropDownContent = document?.querySelector('.dropdown-content');
const overlayer = document.getElementById('overlayer')


burgerButton?.addEventListener('click', () => {
   if(document.querySelector('.hideNavbar')){
        document.querySelector('.hideNavbar').style.animationPlayState = "running";
   }
   navBarContent.classList.toggle('hideNavbar')
   navBarContent.classList.toggle('showNavbar')
   hideScroll()
   overlayer.classList.toggle("overlay")
   if(document.querySelector('.overlay')){
       document.querySelector('.overlay').style.top = navBarContent.offsetHeight + 62 + "px"
   }
})

function hideScroll(){
    if(navBarContent.classList.contains('showNavbar')){
        document.body.style.overflow = 'hidden';
        return
    }
    document.body.style.overflow = 'auto';
}

dropdownMenu?.addEventListener('click', () => {
    if(dropDownContent.style.display === 'block'){
        dropDownContent.style.display = 'none'
    } else if(dropDownContent.style.display === 'none' 
                || dropDownContent.style.display === ''){
        dropDownContent.style.display = 'block'
    }
    document.querySelector('.overlayer').style.top = navBarContent.offsetHeight + 63 + "px"
})