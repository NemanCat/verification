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
                {{statuses[slotProps.data.status].status}}    
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
              <div class="col-12 col-md-2">
                <div class="form-group">
                  <label for="birthday">Дата рождения</label>
                  <InputMask class="p-inputtext-sm form-control" id="birthday" 
                    mask="99.99.9999" v-model.trim="birthday" disabled />
                </div>  
              </div>  
            </div>  
            <div class="form-row">
              <div class="col-12 col-md-3">
                <div class="form-group">
                  <label for="inn">ИНН</label>
                  <InputText :disabled="dataProcess==true" class="p-inputtext-sm form-control" id="inn" 
                    v-model.trim="inn" disabled/>
                </div>
              </div>
              <div class="col-12 col-md-9">
                <div class="form-group">
                  <label for="birth-place">Место рождения</label>
                  <InputText class="p-inputtext-sm form-control" id="birth-place" 
                    v-model.trim="birth_place" disabled />
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
                  <label for="passport">Паспортные данные (серия, номер, кем и когда выдан, код подразделения)</label>
                  <InputText class="p-inputtext-sm form-control" id="passport" 
                    v-model.trim="passport" disabled />
                </div>  
              </div>
            </div>  
            <div class="form-row">          
              <div class="col-12">                
                <div class="form-group">
                  <label for="comment">
                    Комментарий
                  </label>
                  <div class="p-inputgroup">
                    <Textarea class="p-inputtext-sm form-control" id="comment" v-model="comment" /> 
                    <Button icon="pi pi-check" title="Сохранить комментарий" :disabled="!comment"
                      @click="saveComment" />
                  </div>
                </div>  
               
              </div>
            </div>
            <div class="form-row">
              <div class="col-12">
                <div class="form-check form-check-inline">
                  <RadioButton name="isMoscowRadioOptions" id="moscowRadioButton" :value="true" v-model="is_moscow" :disabled="true"/>
                  <label class="form-check-label" for="moscowRadioButton">&nbsp;<b>Москва и МО</b></label>
                </div>
                <div class="form-check form-check-inline">
                  <RadioButton name="isMoscowRadioOptions" id="notMoscowRadioButton" :value="false" v-model="is_moscow" :disabled="true"/>
                  <label class="form-check-label" for="notMoscowRadioButton">&nbsp;<b>Регион</b></label>
                </div>
              </div>
            </div>
            <div class="form-row mt-2">
                <div class="col-6">
                  <fieldset class="border p-2">                      
                    <DataTable ref="phones_dt" :value="this.phones" class="mb-2">            
                      <Column field="phone" header="Номер телефона"></Column>              
                    </DataTable>
                  </fieldset>
                </div>  
                <div class="col-6">
                  <fieldset class="border p-2">                                                       
                      <DataTable ref="emails_dt" :value="this.emails">            
                        <Column field="email" header="Адрес электронной почты"></Column>              
                      </DataTable>
                  </fieldset>
                </div>
              </div>
            <div class="form-row mt-3">
              <div class="col-12">
                <DataView :value="scans" layout="grid" :paginator="false">
                  <template #empty>Нет ни одного скана</template>
                  <template #header>                
                    <div class="p-grid p-nogutter">
                      <div class="p-col-6" style="text-align: left">
                        Сканы экрана БД "Кронос"
                      </div>                  
                    </div>
                  </template>  
                  <template #grid="slotProps">
                    <div class="p-col-12 p-md-4">
                      <Card class="h-100">  
                        <template #header>                          
                          <Button icon="pi pi-search-plus" title="Просмотр скана" 
                            class="p-button-rounded float-left mb-2 ml-2" 
                            @click="current_scan = slotProps.data.img; showScanWindow()" />                                                      
                        </template>
                        <template #content>  
                          <img :src="slotProps.data.img" :alt="slotProps.data.filename" class="center scan-preview"/>    
                        </template>  
                      </Card> 
                    </div>  
                  </template>  
                </DataView>  
              </div>  
            </div>  
            <div class="form-row">
                <div class="col-12">
                  <DataView :value="passport_scans" layout="grid" :paginator="false">
                    <template #empty>Нет ни одного скана</template>
                    <template #header>                           
                      <div class="p-grid p-nogutter">
                        <div class="p-col-6" style="text-align: left">
                          Сканы паспортиста
                        </div>                          
                      </div>
                    </template>  
                    <template #grid="slotProps">
                      <div class="p-col-12 p-md-4">
                        <Card class="h-100">  
                          <template #header>                          
                            <Button icon="pi pi-search-plus" title="Просмотр скана" 
                              class="p-button-rounded float-left mb-2 ml-2"                               
                              @click="current_scan = slotProps.data.img; showScanWindow();" />                                
                          </template>
                          <template #content>  
                            <img :src="slotProps.data.img" :alt="slotProps.data.filename" class="center scan-preview"/>    
                          </template>  
                        </Card> 
                      </div>  
                    </template>  
                  </DataView>  
                </div>  
              </div>
          </form>                    
        </Panel>     
        <Panel>
          <template #header>
            <b>История обработки анкеты</b>
          </template>              
            <div class="row">
              <div class="col-6">Дата и время добавления в БД</div>            
              <div class="col-3">{{inserted}}</div>            
              <div class="col-3">{{operator}}</div>            
            </div>
            <div class="row mt-2">
              <div class="col-6">Дата и время обработки аудитором</div>            
              <div class="col-3">{{audited}}</div>            
              <div class="col-3">{{auditor}}</div>            
            </div>
            <div class="row mt-2">
              <div class="col-6">Дата и время обработки паспортистом</div>            
              <div class="col-3">{{pasported}}</div>            
              <div class="col-3">{{pasportist}}</div>            
            </div>
            <div class="row mt-2">
              <div class="col-6">Дата и время обработки телефонистом</div>            
              <div class="col-3">{{called}}</div>         
              <div class="col-3">{{telephonist}}</div>               
            </div>
        </Panel> 
      </div>        

    </div>

    
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
          {id: 0, status: "не обработана"},
          {id: 1, status: "верифицирована аудитором"},
          {id: 2, status: "забракована аудитором"},
          {id: 3, status: "в обработке у паспортиста"},
          {id: 4, status: "обработана паспортистом"},
          {id: 5, status: "забракована паспортистом"},
          {id: 6, status: "в обработке у телефониста"},
          {id: 7, status: "все телефоны брак"},
          {id: 8, status: "отказался от сотрудничества"},
          {id: 9, status: "готов к сотрудничеству"},
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
        //список сканов экрана оператора
        scans: [],
        //список сканов экрана паспортиста
        passport_scans: [],
        //текущий скан для просмотра
        current_scan: "",
        //список контактных телефонов
        phones: [],
        //список адресов электронной почты
        emails: [],
        //--------------------------------------------------        
        //поля анкеты
        //фамилия
        lastname: "",
        //имя
        firstname: "",
        //отчество
        middle_initial: "",
        //дата рождения
        birthday: "",
        //ИНН
        inn: "",
        //место рождения
        birth_place: "",
        //адрес регистрации
        address: "",
        //паспортные данные
        passport: "",
        //флаг Москва / регионы
        is_moscow: null,
        //комментарий
        comment: "",
        //этапы прохождения анкетой обработки
        inserted: "",
        operator: "",
        audited: "",
        auditor: "",
        pasported: "",
        pasportist: "",
        called: "",
        telephonist: "",
        filters: {},        
        filter_date: "",
        filter_list: [],

        prev_status_filter: null
      }
    },    
    //метаданные страницы
    head() {        
        return {
            //заголовок браузера
            title: 'Анкеты клиентов | АРМ суперпользователя ИС Верификатор'
        }
    },    
    //загрузка статуса приложения
    async beforeCreate() {
       this.$store.commit('initializeStore');
       this.$store.commit('CHANGE_TITLE','Анкеты клиентов');
    },
    //--------------------------------------------------------
    //загрузка страницы    
    async mounted() {            
      this.errorMessageVisible = false;
      this.dataProcess = true;        
      try {
        const res = await this.$nuxt.$http.$post('/api/superuser/get-customers-list');
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
        this.filter_list = this.list;
        if (this.list.length > 0) {          
          this.current_record = this.list[0];                                   
          //загружаем данные первой анкеты
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
        await this.showCustomer();                
      },
      async showScanWindow() {
        const newTab = window.open();
        (<any>newTab).document.body.innerHTML = '<img src="' + this.current_scan + '" style="max-width:1200px;max-height:750px;">';
      },
      //---------------------------------------------------
      //заполнение формы данными текущей выделенной записи
      async showCustomer() {        
        //выбираем данные анкеты
        let searchParams = new URLSearchParams();
        searchParams.set('id',(<any>this.current_record).id.toString());
        this.errorMessageVisible = false;
        try {
          this.dataProcess = true;
          const res = await this.$nuxt.$http.$post('/api/superuser/get-customer',searchParams);
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
        this.birthday = moment(customer.birthday).format("DD.MM.YYYY");
        this.inn = customer.inn;
        this.birth_place = customer.birth_place;
        this.address = customer.address;
        this.passport = customer.passport_data;
        this.comment = customer.auditor_comment == null ? "" : customer.auditor_comment;        
        this.is_moscow = customer.is_moscow;
        this.phones = customer.phones == null ? [] : customer.phones;     
        this.emails = customer.emails == null ? [] : customer.emails;   
        this.scans = customer.screen_scans == null ? [] : customer.screen_scans;
        this.passport_scans = customer.passport_scans == null ? [] : customer.passport_scans;
        this.inserted = this.formatDateTime(customer.inserted);
        this.audited = customer.audited == null ? "" : this.formatDateTime(customer.audited);
        this.pasported = customer.pasported == null ? "" : this.formatDateTime(customer.pasported);
        this.called = customer.called == null ? "" : this.formatDateTime(customer.called);  
        this.operator = customer.operator;
        this.auditor = customer.auditor;
        this.pasportist = customer.pasportist;
        this.telephonist = customer.telephonist;
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
      //сохранение комментария
      async saveComment() {
        this.errorMessageVisible = false;        
        const id = (<any>this.current_record).id;        
        let searchParams = new URLSearchParams();                
        searchParams.set('id',(<any>this.current_record).id.toString());
        searchParams.set('comment',this.comment);        
        try {
          this.dataProcess = true;
          const res = await this.$nuxt.$http.$post('/api/superuser/save-comment',searchParams);
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
        this.showCustomer();
      },
      //--------------------------------------------------
      //фильтрация по статусу анкеты  
      async onFilter(event: any) {                 
        if ((<any>this.filters).status != this.prev_status_filter)   {
          this.prev_status_filter = (<any>this.filters).status;
          if ((<any>this.filters).status == null) {
            //фильтр сброшен
            this.current_record = this.list[0];
          } else {
            //фильтр установлен
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
