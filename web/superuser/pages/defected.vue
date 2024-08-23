<template>
  <div class="container-fluid">    
    <div class="loading" v-show="dataProcess">
         <img src="~assets/img/ajax-loader.gif" alt="Загрузка страницы..." class="block-ui-image"/>
    </div>
    <Message v-if="errorMessageVisible" :severity="this.$store.state.fetchError.fetchErrorLevel==0 ? 'error' : 'warn'" >{{this.$store.state.fetchError.fetchErrorMessage}}</Message>  
    <div class="row">
      <div class="col-6">        
	       <DataTable :value="this.list" class="p-datatable-gridline p-datatable-sm"
              :paginator="true" :rows="10" :filters="filters"
              paginatorTemplate="CurrentPageReport FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink RowsPerPageDropdown"
              :rowsPerPageOptions="[10,20,50]"
              currentPageReportTemplate="Записи от {first} до {last} из {totalRecords}"
              @page="onPageSelect($event)" @filter="onFilter"
              selectionMode="single" :selection="current_record" ref="dt" @row-select="onRowSelect">
            <template #header>
              <div class="table-header">
                Список анкет  
              </div>
            </template>  
            <template #empty>
              Нет ни одной записи
            </template>          
            <Column field="lastname" header="Фамилия" bodyClass="column-text-small" filterField="lastname" filterMatchMode="contains">
              <template #filter>
                <InputText type="text" v-model="filters['lastname']" class="p-column-filter p-inputtext-sm" placeholder="Фамилия"/>
              </template>
            </Column>
            <Column field="firstname" header="Имя" bodyClass="column-text-small"></Column>
            <Column field="middle_initial" header="Отчество" bodyClass="column-text-small"></Column>
            <Column field="status" header="Статус" bodyClass="column-text-small" filterMatchMode="equals" headerStyle="width:30%">
                <template  #body="slotProps">
                    {{slotProps.data.status == 3 ? 'забракована оператором' : (slotProps.data.status == 5 ? 'забракована паспортистом' : 'забракована аудитором')}}    
                </template>
                <template #filter>
                    <Dropdown v-model="filters['status']" :options="statuses" optionLabel="status" 
                        optionValue="id" placeholder="Статус анкеты" class="p-column-filter" :showClear="true">                
                    </Dropdown>
                </template>
            </Column>
          </DataTable> 
          
      </div>
      <div class="col-6">
        <Panel>
          <template #header>            
                <b>Данные анкеты</b>            
          </template>
          <Toolbar>                                
            <template #left>     
              <strong>
              Статус:&nbsp;{{current_record.status == 3 ? 'забракована оператором' : (current_record.status == 5 ? 'забракована паспортистом' : 'забракована аудитором')}}
              </strong>   
            </template>
            <template #right>        
              <Button label="Вернуть в обработку" icon="pi pi-check" class="p-button-success" 
                title="Вернуть анкету в базу данных для повторной обработки"
                @click="returnDialogVisible = true" />          
            </template>
          </Toolbar>
          <form v-if="list.length>0">
            <div class="form-row">
              <div class="col-12 col-md-4">
                <div class="form-group">
                  <label for="lastname">Фамилия</label>
                  <InputText class="p-inputtext-sm form-control" id="lastname" autofocus 
                    v-model.trim="lastname" disabled />
                </div>  
              </div>
              <div class="col-12 col-md-3">
                <div class="form-group">
                  <label for="firstname">Имя</label>
                  <InputText class="p-inputtext-sm form-control" id="firstname" 
                    v-model.trim="firstname" disabled/>
                </div>  
              </div>
              <div class="col-12 col-md-3">
                <div class="form-group">
                  <label for="middle-initial">Отчество</label>
                  <InputText class="p-inputtext-sm form-control" id="middle-initial" 
                    v-model.trim="middle_initial" disabled />
                </div>  
              </div>
            </div>              
            
            <div class="form-row">          
              <div class="col-12">
                <div class="form-group">
                  <label for="address">Адрес регистрации</label>
                  <InputText class="p-inputtext-sm form-control" id="address" 
                    v-model.trim="address" disabled /> 
                </div>  
              </div>
            </div>              
            <div class="form-row">          
              <div class="col-12">                
                <div class="form-group">
                  <label for="comment">
                    Комментарий
                  </label>
                  <Textarea class="p-inputtext-sm form-control" id="comment" v-model="comment" /> 
                </div>  
              </div>
            </div>
          </form>    
          <div class="row">
              <div class="col-2">Оператор:</div>
              <div class="col-10">{{current_record.operator}}</div>
          </div>                
          <div class="row">
              <div class="col-2">Аудитор:</div>
              <div class="col-10">{{current_record.auditor}}</div>
          </div>                
          <div class="row">
              <div class="col-2">Паспортист:</div>
              <div class="col-10">{{current_record.pasportist}}</div>
          </div>                
        </Panel>             
      </div>        

    </div>
    <!-- подтверждение возврата анкеты в обработку -->
      <Dialog :modal="true" header="Подтверждение" :visible.sync="returnDialogVisible" class="p-fluid">
        <div class="confirmation-content">
          <i class="pi pi-exclamation-triangle p-mr-3" style="font-size: 2rem" />&nbsp;
          Вы уверены что хотите вернуть анкету в повторную обработку?
        </div>
        <template #footer>
          <Button label="Нет" icon="pi pi-times" class="p-button-text" @click="returnDialogVisible = false"/>
          <Button label="Да" icon="pi pi-check" class="p-button-text" @click="returnCustomer" />
        </template>
      </Dialog>

    
    </div>
</template>

<script lang="ts">
import Vue from 'vue';
const moment = require('moment');

import Button from 'primevue/button';
import Calendar from 'primevue/calendar';
import Card from 'primevue/card';
import Column from 'primevue/column';
import DataTable from 'primevue/datatable/DataTable';
import DataView from 'primevue/dataview';
import Dialog from 'primevue/dialog';
import Dropdown from 'primevue/dropdown';
import InputMask from 'primevue/inputmask';
import InputText from 'primevue/inputtext';
import Message from 'primevue/message';
import Panel from 'primevue/panel';
import RadioButton from 'primevue/radiobutton';
import TabPanel from 'primevue/tabpanel';
import TabView from 'primevue/tabview';
import Textarea from 'primevue/textarea';
import Toolbar from 'primevue/toolbar';

Vue.component('Button', Button);
Vue.component('Calendar', Calendar);
Vue.component('Card', Card);
Vue.component('Column', Column);
Vue.component('DataTable', DataTable);
Vue.component('DataView', DataView);
Vue.component('Dialog', Dialog);
Vue.component('Dropdown', Dropdown);
Vue.component('InputMask',InputMask);
Vue.component('InputText',InputText);
Vue.component('Message',Message);
Vue.component('Panel',Panel);
Vue.component('RadioButton',RadioButton);
Vue.component('TabPanel',TabPanel);
Vue.component('TabView',TabView);
Vue.component('Textarea', Textarea);
Vue.component('Toolbar', Toolbar);


export default Vue.extend({
  //данные страницы
    data() {
      return { 
        //статусы анкеты
        statuses: [
          {id: 3, status: "забракована оператором"},  
          {id: 2, status: "забракована аудитором"},          
          {id: 5, status: "забракована паспортистом"}
        ],
        //флаг сохранения/загрузки данных
        dataProcess: false,
        //флаг отображения сообщения об ошибке
        errorMessageVisible: false,
        //флаг процесса сохранения данных
        isSubmitting: false,        
        //список неверифицированных анкет
        list: [],
        //текущая запись
        current_record: [],        
        returnDialogVisible: false,
        //--------------------------------------------------        
        //поля анкеты
        //фамилия
        lastname: "",
        //имя
        firstname: "",
        //отчество
        middle_initial: "",        
        //адрес регистрации
        address: "",
        //комментарий
        comment: "",
        filters: {},
        prev_status_filter: null,
        filtered_list: null,
      }
    },    
    //метаданные страницы
    head() {        
        return {
            //заголовок браузера
            title: 'Забракованные анкеты | АРМ суперпользователя ИС Верификатор'
        }
    },    
    //загрузка статуса приложения
    async beforeCreate() {
       this.$store.commit('initializeStore');
       this.$store.commit('CHANGE_TITLE','Забракованные анкеты');
    },
    //--------------------------------------------------------
    //загрузка страницы    
    async mounted() {            
      this.errorMessageVisible = false;
      this.dataProcess = true;        
      try {
        const res = await this.$nuxt.$http.$post('/api/superuser/get-defected-customers-list');
        this.dataProcess = false;        
        if(!res.Success) {                        
          this.$store.commit('CHANGE_FETCHERROR', {
            fetchErrorFlag: true,
            fetchErrorCode: 0,
            fetchErrorMessage: res.Message,
            fetchErrorLevel: res.Level,
          });
          this.errorMessageVisible = true;
          return;
        }        
        
        this.list = JSON.parse(res.Data);     
        if (this.list.length > 0) {          
          this.current_record = this.list[0];        
          await this.showCustomer();
        } 
      }  catch (e) {        
         this.dataProcess = false;         
         if (this.$store.state.fetchError.fetchErrorCode == 403) {
            //отправляем на страницу авторизации
            this.$nuxt.context.app.$cookies.remove('verification-session');
            ( <any>window).document.location.replace('login');
          } else {
            this.errorMessageVisible = true;
          }
          return;
        }
    },   
    //------------------------------------------------------------------------------------
    methods: {
      //формат даты и времени
      formatDateTime(value: string) : string {
        return moment(value).format("DD.MM.YYYY HH:mm");
      },

      //выбор страницы в пейджере
      async onPageSelect(event: any) {                  
        //выбираем первую запись на странице        
        this.current_record = this.list[event.first];
        await this.showCustomer();    
      },

      //выделение строки таблицы
      async onRowSelect(event: any) {
        const id: number = event.data.id;
        for (let i=0; i<this.list.length;i++) {
          if (this.list[i]['id'] === id) {
            this.current_record = this.list[i];
            break;
          }
        }        
        this.showCustomer();                
      },      
      //---------------------------------------------------
      //заполнение формы данными текущей выделенной записи
      async showCustomer() {        
        //выбираем данные анкеты
        let searchParams = new URLSearchParams();        
        searchParams.set('id',(<any>this.current_record).id.toString());
        searchParams.set('source',(<any>this.current_record).source.toString());
        this.errorMessageVisible = false;
        try {
          this.dataProcess = true;
          const res = await this.$nuxt.$http.$post('/api/superuser/get-defected-customer',searchParams);
          this.dataProcess = false;          
          if(!res.Success) {              
            this.$store.commit('CHANGE_FETCHERROR', {
              fetchErrorFlag: true,
              fetchErrorCode: 0,
              fetchErrorMessage: res.Message,
              fetchErrorLevel: res.Level,
            });
            this.errorMessageVisible = true;
            return;
          }
        //заполняем поля анкеты        
        const customer = JSON.parse(res.Data)[0];     
        this.lastname = customer.lastname;        
        this.firstname = customer.firstname;
        this.middle_initial = customer.middle_initial;
        this.address = customer.address;
        this.comment = customer.comment == null ? "" : customer.comment;        
      } catch (e) {        
         this.dataProcess = false;         
         if (this.$store.state.fetchError.fetchErrorCode == 403) {
            //отправляем на страницу авторизации
            this.$nuxt.context.app.$cookies.remove('verification-session');
            ( <any>window).document.location.replace('login');
          } else {
            this.errorMessageVisible = true;
          }
          return;
        }
      },
      //возвращение анкеты в обработку
      async returnCustomer() {
        this.errorMessageVisible = false;
        this.returnDialogVisible = false;
        const id = (<any>this.current_record).id;
        let searchParams = new URLSearchParams();                
        searchParams.set('id',(<any>this.current_record).id.toString());
        searchParams.set('source',(<any>this.current_record).source.toString());
        searchParams.set('status',(<any>this.current_record).status.toString());
        
        try {
          this.dataProcess = true;
          const res = await this.$nuxt.$http.$post('/api/superuser/return-defected-customer',searchParams);
          this.dataProcess = false;     
          if(!res.Success) {              
            this.$store.commit('CHANGE_FETCHERROR', {
              fetchErrorFlag: true,
              fetchErrorCode: 0,
              fetchErrorMessage: res.Message,
              fetchErrorLevel: res.Level,
            });
            this.errorMessageVisible = true;
            return;
          }        
        } catch (e) {                  
           this.dataProcess = false;    
          if (this.$store.state.fetchError.fetchErrorCode == 403) {
            //отправляем на страницу авторизации
            this.$nuxt.context.app.$cookies.remove('verification-session');
            ( <any>window).document.location.replace('login');
          } else {
            this.errorMessageVisible = true;
          }
          return;
        }
        //удаляем из локального списка
        let rec_index = 0;
        if (this.prev_status_filter != null) {
          for (rec_index=0; rec_index<(<any>this.filtered_list).length; rec_index++) {
            if ((<any>this.filtered_list)[rec_index].id == id) {
              this.list = this.list.filter(x => (<any>x).id != id);
              this.filtered_list = (<any>this.filtered_list).filter((x: any) => (<any>x).id != id);
              if (rec_index >= (<any>this.filtered_list).length) rec_index--;
              break;
            }
          }
          if ((<any>this.filtered_list).length > 0) {                          
            this.current_record = (<any>this.filtered_list)[rec_index];
            this.showCustomer();
          }
        } else {
          for (rec_index=0; rec_index<this.list.length; rec_index++) {
            if ((<any>this.list[rec_index]).id == id) {
              this.list.splice(rec_index,1);
              if (rec_index >= this.list.length) rec_index--;
              break;
            }
          }
          if (this.list.length > 0) {                          
            this.current_record = this.list[rec_index];
            this.showCustomer();
          }
        }
                  
            
      },
      //--------------------------------------------------------------
      //фильтрация по статусу анкеты  
      async onFilter(event: any) {                 
        if ((<any>this.filters).status != this.prev_status_filter)   {
          this.prev_status_filter = (<any>this.filters).status;
          if ((<any>this.filters).status == null) {
            //фильтр сброшен
            this.current_record = this.list[0];
            this.filtered_list = null;
          } else {
            //фильтр установлен
            this.filtered_list = event.filteredValue;
            for (let i=0; i<this.list.length;i++) {
              if (this.list[i]['id'] == event.filteredValue[0].id) {
                this.current_record = this.list[i];
                break;
              }
            }
          }
          this.showCustomer(); 
        }
      },

    },  

    
})
</script>
<style>
  .center {
    display: block;
    margin-left: auto;
    margin-right: auto;  
  }

  .scan-preview {
    width: 100%;
    height: 100%;
  }

   
</style>
