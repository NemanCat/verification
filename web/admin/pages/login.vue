<template>
<client-only>
    <div class="container-fluid h-100">
        <div class="row justify-content-center align-items-center h-100">
            <div class="col col-sm-7 col-md-7 col-lg-5 col-xl-4 pt-3">
                <div class="card">
                    <div class="card-header text-center"><h4>Авторизация в системе</h4></div>
                        <div class="card-body">        
                            <Message v-if="this.$store.state.fetchError.fetchErrorFlag==true" severity="error" ><small>{{this.$store.state.fetchError.fetchErrorMessage}}</small></Message>                
                            <div class="form-group row">                                
                                <label for="email" class="col-sm-3 col-form-label">Имя входа</label>
                                <div class="col-sm-9">
                                    <InputText id="login" required="true" placeholder="Логин" style="width:100%" v-model.trim="login" :class="{'p-invalid': isSubmitted && !login}"/>
                                    <small v-if="isSubmitted && !login" class="p-invalid">Укажите имя входа в систему!</small>           
                                </div>                              
                            </div>
                            <div class="form-group row">
                                <label for="password" class="col-sm-3 col-form-label">Пароль</label>
                                <div class="col-sm-9">
                                    <Password id="password" required="true" placeholder="Пароль" style="width:100%" v-model.trim="password" :feedback="false" :class="{'p-invalid': isSubmitted && !password}"/>
                                    <small v-if="isSubmitted && !password" class="p-invalid">Укажите пароль!</small>
                                </div>
                               
                            </div>
                            <div class="form-group">
                                <button class="btn btn-info btn-lg btn-block" @click="loginClick">Войти</button>
                            </div>
                       
                    </div>
                </div>
            </div>
        </div>
    </div>
</client-only>    
</template>

<script lang="ts">
import Vue from 'vue'

export default Vue.extend({   
    layout: "fullscreen",
    //данные страницы
    data() {
        return {         
            isSubmitted: false,
            login: "",
            password: "",
        }
    },
    //метаданные страницы
    head() {        
        return {
            //заголовок браузера
            title: 'Авторизация в системе | АРМ администратора ИС Верификатор'
        }
    },    
    //загрузка статуса приложения
    async beforeCreate() {
       this.$store.commit('initializeStore');
       this.$store.commit('CHANGE_TITLE','Авторизация в системе');
    },
    methods: {
        //авторизация
        async loginClick() {            
            this.$store.commit('CHANGE_FETCHERROR', {
                fetchErrorFlag: false,
                fetchErrorCode: 0,
                fetchErrorMessage: ''
            });
            this.isSubmitted = true;
            if (this.login && this.password) {                
                //логинимся
                (<any>window).document.getElementById('loading').style = 'display: block;';
                this.$nuxt.$http.setHeader('Content-Type', 'application/x-www-form-urlencoded');
                const searchParams = new URLSearchParams();
                searchParams.set('login',this.login);
                searchParams.set('password',this.password);   
                searchParams.set('category','0');
                let res: Response = new Response; 
                try {                       
                    res = await this.$nuxt.$http.$post(`/api/authorization`,searchParams);                                           
                    this.isSubmitted = false;
                    //авторизация пройдена, создаём cookie и переходим на начальную страницу АРМа
                    const data = JSON.parse(JSON.stringify(res));
                    this.$nuxt.$cookies.set('verification-session',data.sessionid);
                    this.$nuxt.$store.commit('CHANGE_CSFR_SECURITY_TOKEN',data.csfr_token);
                    this.$nuxt.$store.commit('CHANGE_USERID',data.id);
                    this.$nuxt.$store.commit('CHANGE_USERNAME',data.name);
                    (<any>window).document.location.replace('');
                } catch (e) {       
                    (<any>window).document.getElementById('loading').style = 'display: none;';
                    this.isSubmitted = false;
                }
             }         
        },
    }
    
    //--------------------------------------------------------
   
     


})
</script>