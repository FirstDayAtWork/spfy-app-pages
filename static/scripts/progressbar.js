"use strict"

const loader = document?.querySelector('#loader');

function animate({duration, draw, timing}) {
    let start = performance.now();
  
    requestAnimationFrame(function animate(time) {
      let timeFraction = (time - start) / duration;
      if (timeFraction > 1) timeFraction = 1;
  
      let progress = timing(timeFraction)
  
      draw(progress);
  
      if (timeFraction < 1) {
        requestAnimationFrame(animate);
      }
  
    });
  }

    
  document.addEventListener('readystatechange', () => {
    if (document.readyState !== "complete") {
      loader.style.visibility = 'visible'
          animate({
              duration: 10000,
              timing: function(timeFraction) {
                return timeFraction;
              },
              draw: function(progress) {
                  loader.value = progress * 100;
              }
            });
    } 
        loader.style.visibility = 'hidden'
        return
}) 
