<template>
    <div>        
        <client-only>
            <Message v-if="this.fetchErrorPageFlag==true" severity="error" >{{this.$store.state.fetchError.fetchErrorMessage}}</Message>
            <DataTable :value="this.list" class="p-datatable-striped p-datatable-gridline p-datatable-sm" autoLayout 
                selectionMode="single"  ref="dt" :loading="this.dataLoading" @row-select="onRowSelect">
                <template #header>
                    <div class="table-header">
                        <Button label="Добавить" title="Добавление новой записи" icon="pi pi-plus" class="p-button-success" @click="newRecord" />
                        <Button icon="pi pi-refresh" title="Обновить данные" @click="refresh" />
                    </div>
                </template>
                <template #empty>
                    Нет ни одной записи
                </template>
                <Column field="name" header="ФИО" bodyClass="column-text-small"></Column>
                <Column field="login" header="Логин" bodyClass="column-text-small"></Column>
                <Column field="category" header="Категория" bodyClass="column-text-small">
                  <template #body="slotProps">
                    {{categories[slotProps.data.category].name}}
                  </template>
                </Column>
                <Column bodyClass="column-text-small">
                    <template #body="slotProps">
                        <div style="float:right;">
                            <Button title="Редактировать запись" icon="pi pi-pencil" class="p-button-rounded p-button-success p-mr-2" @click="editRecord(slotProps.data)" />
                            &nbsp;
                            <Button title="Удалить запись" icon="pi pi-trash" class="p-button-rounded p-button-danger" @click="confirmDeleteRecord(slotProps.data)" />
                        </div>  
                    </template>
                </Column>
            </DataTable>  
            <Dialog v-if="list.length>0" :modal="true" header="Подтверждение" :visible.sync="deleteDialogVisible" class="p-fluid">
      <div class="confirmation-content">
        <i class="pi pi-exclamation-triangle p-mr-3" style="font-size: 2rem" />&nbsp;
        <span v-if="current_record">Вы уверены что хотите удалить пользователя <b>{{current_record.name}}</b>?</span>
    </div>
    <template #footer>
        <Button label="Нет" icon="pi pi-times" class="p-button-text" @click="deleteDialogVisible = false"/>
        <Button label="Да" icon="pi pi-check" class="p-button-text" @click="deleteRecord" />
    </template>
    </Dialog> 

     <Dialog :closable="savingData==false"  :style="{width: '450px'}" :visible.sync="editDialogVisible" :header="editDialogTitle" :modal="true" class="p-fluid">       
      <div id="saving_data" v-if="savingData==true"><img id="loading-image" src="~assets/img/ajax-loader.gif" alt="Загрузка страницы..." /></div>
      <Message v-if="this.fetchErrorDialogFlag==true" severity="error" >{{this.$store.state.fetchError.fetchErrorMessage}}</Message> 
      <div class="p-field p-mb-3">
        <label for="name">ФИО</label>
        <InputText :disabled="savingData==true" class="p-inputtext-sm" id="name" autofocus required="true" 
          v-model.trim="current_record.name"
          :class="{'p-invalid': isSubmitted && !current_record.name}" />
        <small class="p-invalid" v-if="isSubmitted && !current_record.name">Укажите ФИО пользователя!</small>       
      </div>
      <div class="p-field">
        <label for="category">Категория пользователя</label>
        <Dropdown v-model="current_record.category" :options="categories" optionLabel="name" optionValue="category" placeholder="Выберите категорию пользователя" />
      </div>
      <div class="p-field">
        <label for="login">Логин</label>
        <InputText :disabled="savingData==true" class="p-inputtext-sm" id="login" required="true" 
          v-model.trim="current_record.login"
          :class="{'p-invalid': isSubmitted && !current_record.login}"/>
        <small class="p-invalid" v-if="isSubmitted && !current_record.login">Укажите имя входа в систему!</small>       
      </div>
      <div class="p-field">
        <label for="password">Пароль</label>
        <Password autocomplete="new-password" :disabled="savingData==true" class="p-inputtext-sm" id="password" v-model.trim="current_record.password" required="true" 
          :class="{'p-invalid': isSubmitted && !current_record.password}" 
          promptLabel="Введите пароль" weakLabel="Слабый пароль"
          mediumLabel="Средний пароль" strongLabel="Сильный пароль"/>
        <small class="p-invalid" v-if="isSubmitted && !current_record.password">Укажите пароль пользователя!</small>       
        <small class="p-invalid" v-if="isSubmitted && current_record.confirm_password && current_record.password && current_record.confirm_password != current_record.password">Пароли не совпадают!</small>
      </div>
      <div class="p-field">
        <label for="confirm_password">Подтверждение пароля</label>
        <Password autocomplete="new-password" :disabled="savingData==true" class="p-inputtext-sm" id="confirm_password" v-model.trim="current_record.confirm_password" required="true" 
          :class="{'p-invalid': isSubmitted && !current_record.confirm_password}" 
          :feedback=false />
        <small class="p-invalid" v-if="isSubmitted && !current_record.confirm_password">Подтвердите пароль пользователя!</small>       
        <small class="p-invalid" v-if="isSubmitted && current_record.confirm_password && current_record.password && current_record.confirm_password != current_record.password">Пароли не совпадают!</small>
      </div>
      <template #footer>
        <Button :disabled="savingData==true" label="Отмена" icon="pi pi-times" class="p-button-text p-button-danger" @click="editDialogVisible = false"/>
        <Button :disabled="savingData==true" label="Сохранить" icon="pi pi-check"  @click="saveRecord" />
      </template>
      
    </Dialog>
        </client-only>
        </div>   
</template>

<script lang="ts">
import Vue from 'vue'
import { DataTable } from 'primevue/datatable';
import { Retvalue } from '~/model/retvalue';

export default Vue.extend({    
    //данные страницы
    data() {
        return {           
            //список категорий пользователей
            categories: [
			        {name: 'Администратор', category: 0},
			        {name: 'Архивариус', category: 1},
			        {name: 'Оператор БД', category: 2},
			        {name: 'Аудитор', category: 3},
              {name: 'Телефонист', category: 4},
              {name: 'Паспортист', category: 5},
              {name: 'Суперпользователь', category: 6},
		        ], 
            list:[],
            table: <DataTable>{},
            //флаг отображения сообщения об ошибке на странице
            fetchErrorPageFlag: false,
            //флаг отображения сообщения об ошибке в диалоге
            fetchErrorDialogFlag: false,
            //флаг видимости подтверждения удаления записи
            deleteDialogVisible: false,
            //флаг видимости окна добавления/редактирования записи
            editDialogVisible: false,
            //заголовок окна добавления/редактирования записи
            editDialogTitle:  '',   
            //флаг отправки данных из формы редактирования
            isSubmitted:  false,
            //текущая выделенная запись
            selectedRecord: null, 
            //текущая выбранная запись
            current_record:  {
                "id": 0,
                "name": "",
                "login": "",
                "password": "",
                "confirm_password": "",
                "category": 0
            },
            //захэшированный пароль - запомнить на время редактирования данных админа
            current_password:  "",
            //сгенерированный случайный пароль
            random_password:  "",
            //флаг загрузки данных
            dataLoading: false,
            //флаг сохранения данных
            savingData: false,
        }
    },
    //метаданные страницы
    head() {        
        return {
            //заголовок браузера
            title: 'Администраторы системы | АРМ администратора ИС Верификатор'
        }
    },    
    //загрузка статуса приложения
    async beforeCreate() {
       this.$store.commit('initializeStore');
       this.$store.commit('CHANGE_TITLE','Администраторы системы');
    },       
    //--------------------------------------------------------
    async mounted() {          
        //получение данных с сервера
        try {          
          (<any>window).document.getElementById('loading').style = 'display: block;';
          this.dataLoading = true;           
          const res = await this.$nuxt.$http.$post('/api/admin/users');
          this.dataLoading = false;
          (<any>window).document.getElementById('loading').style = 'display: none;';
            if(!res.Success) {
                this.$store.commit('CHANGE_FETCHERROR', {
                  fetchErrorFlag: true,
                  fetchErrorCode: 0,
                  fetchErrorMessage: 'При получении данных с сервера произошла ошибка. Сообщение об ошибке: ' + res.Message
                });
                this.fetchErrorPageFlag = true;
                return;
            }
            this.list = JSON.parse(res.Data);              
        } catch (e) {                
          this.dataLoading = false;  
          if (this.$store.state.fetchError.fetchErrorCode == 403) {
            //отправляем на страницу авторизации
            this.$nuxt.context.app.$cookies.remove('verification-session');
            (<any>window).document.location.replace('login');
          } else {
            this.fetchErrorPageFlag = true;
            (<any>window).document.getElementById('loading').style = 'display: none;';
          }
          return;
        } 
        this.table = this.$refs.dt as DataTable;    
        if (this.list.length > 0) this.current_record = this.list[0];
        else this.current_record = {
          "id": 0,
          "name": "",
          "login": "",
          "password": "",
          "confirm_password": "",
          "category": 0
        };
        this.table.selection = this.current_record;      
    },
    methods: {
        //получение данных с сервера
        async refresh() {
          this.fetchErrorPageFlag = false;
          try {          
            (<any>window).document.getElementById('loading').style = 'display: block;';
            this.dataLoading = true;           
            const res = await this.$nuxt.$http.$post('/api/admin/users');
            this.dataLoading = false;
            (<any>window).document.getElementById('loading').style = 'display: none;';
            if(!res.Success) {              
                this.$store.commit('CHANGE_FETCHERROR', {
                  fetchErrorFlag: true,
                  fetchErrorCode: 0,
                  fetchErrorMessage: 'При получении данных с сервера произошла ошибка. Сообщение об ошибке: ' + res.Message
                });
                this.fetchErrorPageFlag = true;
                return;
            }
            this.list = JSON.parse(res.Data);              
          } catch (e) {                
              this.dataLoading = false;  
              if (this.$store.state.fetchError.fetchErrorCode == 403) {
                //отправляем на страницу авторизации
                this.$nuxt.context.app.$cookies.remove('verification-session');
              ( <any>window).document.location.replace('login');
              } else {
                this.fetchErrorPageFlag = true;
                (<any>window).document.getElementById('loading').style = 'display: none;';
              }
              return;
          }  
          if (this.list.length > 0) this.current_record = this.list[0];
          else this.current_record = {
            "id": 0,
            "name": "",
            "login": "",
            "password": "",
            "confirm_password": "",
            "category": 0
          };
          this.table.selection = this.current_record;
        }, 

        //добавление новой записи
         newRecord() {
            this.fetchErrorDialogFlag = false;
            this.isSubmitted = false;
            this.current_record = {
                "id": 0,
                "name": "",
                "login": "",
                "password": "",
                "confirm_password": "",
                "category": 0
            }
            this.editDialogTitle = "Добавление записи";           
            this.editDialogVisible = true;                  
        },
        //редактирование записи
       editRecord(admin: any) {         
            this.fetchErrorDialogFlag = false;
            this.current_record = {...admin};
            //генерируем случайный пароль для заполнения поля пароля
            this.current_password = this.current_record.password;
            this.random_password = Math.random().toString(36).substring(2,15);
            this.current_record.password = this.random_password;
            this.current_record.confirm_password = this.random_password;
            this.isSubmitted = false;
            this.editDialogTitle = "Редактирование записи";            
            this.table.selection = this.current_record;
            this.editDialogVisible = true;
        },

        //сохранение отредактированной записи
      async saveRecord() {       
        this.fetchErrorDialogFlag = false;
        this.isSubmitted = true;
        if (this.current_record.name.trim() && this.current_record.login.trim()
          && this.current_record.password.trim() && this.current_record.confirm_password.trim()
          && this.current_record.password == this.current_record.confirm_password) {       
            let res: Retvalue;     
            const searchParams = new URLSearchParams();
            searchParams.set('id',this.current_record.id.toString());
            searchParams.set('name',this.current_record.name);
            searchParams.set('category',this.current_record.category.toString());
            searchParams.set('login',this.current_record.login);
            searchParams.set('password',this.current_record.password);
            //для существующей записи определяем, изменился ли пароль
            if (this.current_record.id !== 0) {
              if (this.current_record.password != this.random_password) {
                searchParams.set('password_changed','true');
              } else {
                searchParams.set('password_changed','false');
              }
            }
            //сохраняем данные на сервере
            try {                          
                this.savingData = true;                     
                const res = await this.$nuxt.$http.$post('/api/admin/users/edit',searchParams);
                this.savingData = false;
                         
                if(!res.Success) {
                    this.$store.commit('CHANGE_FETCHERROR', {
                        fetchErrorFlag: true,
                        fetchErrorCode: 0,
                        fetchErrorMessage: 'При сохранении данных на сервере произошла ошибка. Сообщение об ошибке: ' + res.Message
                    });
                  this.fetchErrorDialogFlag = true;  
                  return;
                }    
                //изменения успешно сохранены на сервере
                if (this.current_record.id !== 0) {
                  //сохранение изменений существующей записи в локальном списке
                  for (let i=0; i< this.list.length; i++) {
                    if ((<any>this.list[i]).id === this.current_record.id) {
                      (<any>this.list[i]).name = this.current_record.name;
                      (<any>this.list[i]).category = this.current_record.category;
                      (<any>this.list[i]).login = this.current_record.login;
                      (<any>this.list[i]).password = this.current_record.password;
                      (<any>this.list[i]).confirm_password = this.current_record.confirm_password;
                      break;
                    }
                  }
                } else {
                  //добавление новой записи в локальный список
                  const  new_id: number = res.Data as number;
                  (<any>this.list).push({
                      id: new_id,
                      name: this.current_record.name,
                      category: this.current_record.category,
                      login: this.current_record.login,
                      password: this.current_record.password,
                      confirm_password: this.current_record.password
                  });
                  for (let i=0; i<this.list.length;i++) {
                    if ((<any>this.list[i]).id === new_id) {
                      this.current_record = this.list[i];
                      break;
                    }
                  }
                }
                this.table.selection = this.current_record;
                this.list = this.list.sort(function(obj1: any, obj2: any) {return obj1.name.localeCompare(obj2.name)});
                this.editDialogVisible = false; 
            } catch (e) {                
              this.savingData = false;  
              if (this.$store.state.fetchError.fetchErrorCode == 403) {
                //отправляем на страницу авторизации
                this.$nuxt.context.app.$cookies.remove('verification-session');
                ( <any>window).document.location.replace('login');
              } else {
                this.fetchErrorDialogFlag = true;          
              }
              return;
           }
      }   
    },
        //подтверждение удаления записи
        confirmDeleteRecord(user: any) {
            this.current_record = {...user};
            this.table.selection = this.current_record;
            this.deleteDialogVisible = true;      
        },
        //удаление записи
        async deleteRecord() {
            let res: Retvalue;     
            const searchParams = new URLSearchParams();
            searchParams.set('id',this.current_record.id.toString()); 
            try {                          
                this.savingData = true;                                     
                const res = await this.$nuxt.$http.$post('/api/admin/users/remove',searchParams);
                this.savingData = false;                         
                if(!res.Success) {
                    this.$store.commit('CHANGE_FETCHERROR', {
                        fetchErrorFlag: true,
                        fetchErrorCode: 0,
                        fetchErrorMessage: 'При изменении данных на сервере произошла ошибка. Сообщение об ошибке: ' + res.Message
                    });
                  this.fetchErrorDialogFlag = true;  
                  return;
                }
                //удаление записи в локальном списке записей
                for (let i=0; i< this.list.length; i++) {
                  if ((<any>this.list[i]).id === this.current_record.id) {
                    this.list.splice(i--,1);
                    this.current_record = this.list[0];
                    this.table.selection = this.current_record;
                    break;
                  }
                }
                this.deleteDialogVisible = false;
            }  catch (e) {                
              this.savingData = false;  
              if (this.$store.state.fetchError.fetchErrorCode == 403) {
                //отправляем на страницу авторизации
                this.$nuxt.context.app.$cookies.remove('verification-session');
                ( <any>window).document.location.replace('login');
              } else {
                this.fetchErrorDialogFlag = true;          
              }
              return;
           }  
    },
    
    //------------------------------------------
    onRowSelect(event: any) {
        const id: number = event.data.id;
        for (let i=0; i<this.list.length;i++) {
          if (this.list[i]['id'] === id) {
            this.current_record = this.list[i];
            break;
          }
        }
        this.table.selection = this.current_record;      
    }
  }
})
</script>


<style scoped>
 .table-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }
</style>
