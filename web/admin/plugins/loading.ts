(<any>window).onNuxtReady((nuxt: any) => {  
  if(nuxt.$store.state.sidebar_collapsed) {
    (<any>window).document.getElementById('template').className = 'collapsed'  
  } else {    
    (<any>window).document.getElementById('template').className = 'wide'
  }
  (<any>window).document.getElementById('loading').style = 'display: none;'
});
 