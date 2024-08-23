<template>
    <div id="template" ref="template"> 
      <div id="loading"><img id="loading-image" src="~assets/img/ajax-loader.gif" alt="Загрузка страницы..." /></div>
      <!-- menu sidebar -->
      <div class="sidebar">     
        <ul class="list-unstyled">  
          <li v-bind:class="{active: $nuxt.$route.path == '/'}" title="Стартовая страница АРМ Администратора"><a href=""><i class="fa fa-table"></i><span>&nbsp;&nbsp;Панель управления</span></a></li>      
          <li v-bind:class="{active: $nuxt.$route.path == '/users/'}" title="Список пользователей системы"><a href="users"><i class="fa fa-users-cog"></i><span>&nbsp;&nbsp;Пользователи</span></a></li>          
          <hr/>
          <li title="Выйти из системы"><a href="#" @click="logout"><i class="fa fa-door-open"></i><span>&nbsp;&nbsp;Выход</span></a></li>           
        </ul>
      </div>  
      <!-- top navbar -->
      <nav class="navbar sticky-top top-navbar">
        <button class="navbar-toggler" type="button" 
          data-toggle="collapse" data-target="#navbarToggleExternalContent" 
          aria-controls="navbarToggleExternalContent" aria-expanded="false" aria-label="Toggle navigation"
          @click="toggleSidebar">
          <span  v-bind:title="this.$store.state.sidebar_collapsed ? 'Развернуть меню' : 'Свернуть меню'">
            <i id="toggler-icon" v-bind:class="this.$store.state.sidebar_collapsed ? 'fa fa-grip-lines-vertical' : 'fa fa-bars'" style="color:#fff; font-size:16px;"></i>
          </span>
        </button>
        <span class="navbar-brand mb-0 h1 ml-2">АРМ Администратора</span>      
        <div class="btn-group ml-auto">       
           <button type="button" class="btn btn-link shadow-none"  @click="logout" title="Выйти из системы">
            <i class="fa fa-door-open"></i>&nbsp;&nbsp;Выход
          </button>
        </div>  
      </nav>
      <!-- content area -->
      <div class="main">  
        <!-- component content -->
        <div style="padding-bottom:15px;">
            <h4>{{ this.$store.state.title }}</h4>
        </div>
        <div class="pl-2 pr-2 pb-3"><Nuxt/></div>     
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
 
  methods: {    
    async toggleSidebar(event: any) {      
      this.$store.commit('CHANGE_SIDEBAR_COLLAPSED',!this.$store.state.sidebar_collapsed);
      if(this.$store.state.sidebar_collapsed) {
        (<any>this.$refs['template']).className = 'collapsed'  
      } else {    
        (<any>this.$refs['template']).className = 'wide'
      }
    },
    //выход из системы
    async logout(event: any) {
      event.preventDefault();      
      (<any>window).document.getElementById('loading').style = 'display: block;';
     // this.$nuxt.context.app.$http.setHeader('Content-Type', 'application/x-www-form-urlencoded');
     // this.$nuxt.context.app.$http.setHeader('x-request-security-token', this.$store.state.csfr_security_token);
      await this.$nuxt.context.app.$http.$post(`/api/logout`);
      this.$nuxt.context.app.$cookies.remove('verification-session');
      (<any>window).document.location.replace('login');
    }
  },
    
})
</script>