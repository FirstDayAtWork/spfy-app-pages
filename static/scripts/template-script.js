"use strict"

const burgerButton = document.querySelector('.burger-btn');
const navBarContent = document.querySelector('.navbar-content');
const dropdownMenu = document?.querySelector('.title-plus-arrow');
const dropDownContent = document?.querySelector('.dropdown-content');


burgerButton?.addEventListener('click', () => {
   if(document.querySelector('.hideNavbar')){
        document.querySelector('.hideNavbar').style.animationPlayState = "running";
   }
   navBarContent.classList.toggle('hideNavbar')
   navBarContent.classList.toggle('showNavbar')
   hideScroll()
})

function hideScroll(){
    if(navBarContent.classList.contains('showNavbar')){
        document.body.style.overflow = 'hidden';
        return
    }
    document.body.style.overflow = 'auto';
}

dropdownMenu?.addEventListener('click', () => {
    if(dropDownContent.style.visibility === 'visible'){
        dropDownContent.style.visibility = 'hidden'   
    } else if(dropDownContent.style.visibility === 'hidden' 
                || dropDownContent.style.visibility === ''){
                dropDownContent.style.visibility = 'visible'
    }
})