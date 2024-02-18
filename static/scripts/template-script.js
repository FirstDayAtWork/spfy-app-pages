const burgerButton = document.querySelector('.burger-btn');
const navBarContent = document.querySelector('.navbar-content');
const navBarLinks = document.querySelectorAll('.navbar-links');
const content = document.querySelector('.container');
const dropdownMenu = document?.querySelector('.dropdown');
const dropDownContent = document?.querySelector('.dropdown-content');

// show/hide navbar content
burgerButton?.addEventListener('click', () => {
    if(navBarContent.style.display === 'flex'){
        navBarContent.style.display = 'none'
        navBarLinks.forEach(el => {
            el.style.display = 'none'
        })
    } else if(navBarContent.style.display === 'none' 
                || navBarContent.style.display === ''){
        navBarContent.style.display = 'flex'
        navBarLinks.forEach(el => {
            el.style.display = 'flex'
        })
    }
    
})


dropdownMenu?.addEventListener('click', () => {
    if(dropDownContent.style.display === 'block'){
        dropDownContent.style.display = 'none'
    } else if(dropDownContent.style.display === 'none' 
                || dropDownContent.style.display === ''){
        dropDownContent.style.display = 'block'
    }

})


// dropdownMenu?.addEventListener('mouseover', () => {
//         dropDownContent.style.display = 'block'
// })


// dropdownMenu?.addEventListener('mouseout', () => {
//     dropDownContent.style.display = 'none'
// })


// burgerButton.addEventListener('blur', () => {
//     navBarContent.style.display = 'none'
// })

content.addEventListener('click', () => {
    navBarContent.style.display = 'none'
})