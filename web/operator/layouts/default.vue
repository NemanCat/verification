<template>
  <div>
    <BlockUI :blocked="dataProcess" :fullScreen="true" class="block-ui">
      <img id="loading-image" src="~assets/img/ajax-loader.gif" alt="Загрузка страницы..." class="block-ui-image" v-if="dataProcess"/>
    </BlockUI>
    <!-- top navbar -->
    <nav class="navbar navbar-expand-lg navbar-dark" style="background-color:  #535C69">
      <a class="navbar-brand">АРМ оператора</a>
      <button class="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarNav" aria-controls="navbarNavAltMarkup" aria-expanded="false" aria-label="Toggle navigation">
        <span class="navbar-toggler-icon"></span>
      </button>
      <div class="collapse navbar-collapse" id="navbarNav">        
        <div class="navbar-nav justify-content-end ml-auto">          
          <a class="nav-item nav-link" :class="$nuxt.$route.path == '/' ?  'disabled' : 'active'" href=""
            title="Ввод новых анкет клиентов">
            <i class="fas fa-edit"></i>&nbsp;&nbsp;Ввод данных
          </a>          
          <!--
          <a  class="nav-item nav-link" :class="$nuxt.$route.path == '/report/' ?  'disabled' : 'active'" href="report"
            title="Отчёт по введённым анкетам">
            <i class="fas fa-list"></i>&nbsp;&nbsp;Отчёт
          </a>
          -->          
          <a class="nav-item nav-link active" href="#" title="Выход из системы" @click="logout">
            <i class="fa fa-door-open"></i>&nbsp;&nbsp;Выход
          </a>
        </div>
      </div>
    </nav> 
    <div class="main">
      <!-- component content -->
      <div>
          <h4>{{ this.$store.state.title }}</h4>
      </div>
      <ProgressBar v-if="seconds_to_exit <= 100" :value="seconds_to_exit" class="mb-2">
        До выхода из системы {{seconds_to_exit}} сек.
      </ProgressBar>
      <div style="padding-left:10px;padding-right:10px;"><Nuxt/></div>     
      <!-- page footer -->
      <footer class="site-footer">     
        <div class="container">       
          <p class="copyright-text">Информационная система Верификатор</p>      
        </div>
      </footer>
    </div>
  </div>  
</template>

<script lang="ts">
import Vue from 'vue'

export default Vue.extend({
  middleware: 'authenticated',
  data() {
    return {
      dataProcess: false,
      seconds_to_exit: 60*60,
      timerId: null,      
    }
  },
  mounted() {
    //запускаем таймер
    this.seconds_to_exit = 60*60;
    (<any>this.timerId) = setInterval(() => {
      if (this.seconds_to_exit == 0) {
        this.system_logout();
      } else {        
        this.seconds_to_exit--;                
      }      
      }, 1000);
      //event listener on mouse click
      window.addEventListener('mousemove',this.onMouseMoveListener);
      //event listener on keyboard event
      window.addEventListener('keydown',this.onKeyboardListener);
  },
  //сброс listeners при выходе
  beforeDestroy() {
    window.removeEventListener('mousemove',this.onMouseMoveListener);
    window.removeEventListener('keydown',this.onKeyboardListener);
  },
  methods: {        
    //штатный выход из системы
    async logout(event: any) {
      event.preventDefault(); 
      window.removeEventListener('mousemove',this.onMouseMoveListener);
      window.removeEventListener('keydown',this.onKeyboardListener);
      this.dataProcess = true;
      await this.$nuxt.context.app.$http.$post(`/api/logout`);
      this.dataProcess = false;
      this.$nuxt.context.app.$cookies.remove('verification-session');
      (<any>window).document.location.replace('login');
    },
    //выход из системы по таймауту
    async system_logout() {
      window.removeEventListener('mousemove',this.onMouseMoveListener);
      window.removeEventListener('keydown',this.onKeyboardListener);
      this.dataProcess = true;
      await this.$nuxt.context.app.$http.$post(`/api/system-logout`);
      this.dataProcess = false;
      this.$nuxt.context.app.$cookies.remove('verification-session');
      (<any>window).document.location.replace('login');
    },
    //mouse move event listener
    async onMouseMoveListener() {
      this.seconds_to_exit = 60*60;
    },
    //keyboard event listener
    async onKeyboardListener() {
      this.seconds_to_exit = 60*60;
    },
  },
    
})
</script>