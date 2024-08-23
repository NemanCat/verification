<template>
  <div class="container-fluid">    
    <BlockUI :blocked="dataProcess" :fullScreen="true" class="block-ui">
        <img id="loading-image" src="~assets/img/ajax-loader.gif" alt="Загрузка страницы..." class="block-ui-image" v-if="dataProcess"/>
    </BlockUI>
    <Message v-if="errorMessageVisible" :severity="this.$store.state.fetchError.fetchErrorLevel==0 ? 'error' : 'warn'" >{{this.$store.state.fetchError.fetchErrorMessage}}</Message>  
    <div class="row">
      <div class="col-6">        
	       <DataTable :value="this.list" class="p-datatable-gridline p-datatable-sm"
              :paginator="true" :rows="10" :filters="filters"
              paginatorTemplate="CurrentPageReport FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink RowsPerPageDropdown"
              :rowsPerPageOptions="[10,20,50]"
              currentPageReportTemplate="Записи от {first} до {last} из {totalRecords}"
              @page="onPageSelect($event)"
              selectionMode="single" ref="dt" @row-select="onRowSelect">
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
                <InputText type="text" v-model="filters['lastname']" class="p-column-filter p-inputtext-sm" size="10" placeholder="Фамилия"/>
              </template>
            </Column>
            <Column field="firstname" header="Имя" bodyClass="column-text-small"></Column>
            <Column field="middle_initial" header="Отчество" bodyClass="column-text-small"></Column>
            
            <Column field="name" header="Оператор" bodyClass="column-text-small" filterField="name" 
              filterMatchMode="contains">
              <template #filter>
                <InputText type="text" v-model="filters['name']" class="p-column-filter p-inputtext-sm" size="8"  placeholder="Оператор"/>
              </template>
            </Column>
            <Column field="login" header="Логин" bodyClass="column-text-small"></Column>
            <Column header="Дата и время добавления" bodyClass="column-text-small">
              <template #body="slotProps">
                {{formatDateTime(slotProps.data.inserted)}}
              </template>  
              <template #filter>
                <Calendar v-model="filter_date" :locale="ru" :showIcon="false" 
                  dateFormat="dd.mm.yy" class="p-column-filter" appendTo="body"
                  :showButtonBar="true" placeholder="Дата" @date-select="filterDate"
                  @clear-click="filterDate"/>
              </template>
            </Column>      
          </DataTable> 
          
      </div>
      <div class="col-6">
        <Panel>
          <template #header>
            <b>Данные анкеты</b>
          </template>
          <Toolbar v-if="list.length>0">        
            <template #left>
              <Button label="Забраковать" icon="pi pi-times" class="p-button-danger" 
                title="Забраковать анкету"
                @click="defectDialogVisible = true" />
            </template>        
            <template #right>        
              <Button label="Одобрить" icon="pi pi-check" class="p-button-success p-mr-2" 
                title="Сохранить введённые данные и пометить анкету как верифицированную" 
                @click="confirmDialogVisible=true;" />          
            </template>
          </Toolbar>
          <form v-if="list.length>0">
            <div class="form-row">
              <div class="col-12 col-md-4">
                <div class="form-group">
                  <label for="lastname">Фамилия</label>
                  <InputText :disabled="dataProcess==true" class="p-inputtext-sm form-control" id="lastname" autofocus 
                    required="true" maxlength="50" v-model.trim="lastname"
                    :class="{'p-invalid': isSubmitting && !lastname}" />
                  <small class="p-invalid" v-if="isSubmitting && !lastname">Укажите фамилию клиента!</small>                  
                </div>  
              </div>
              <div class="col-12 col-md-3">
                <div class="form-group">
                  <label for="firstname">Имя</label>
                  <InputText :disabled="dataProcess==true" class="p-inputtext-sm form-control" id="firstname" 
                    required="true" maxlength="50" v-model.trim="firstname"
                    :class="{'p-invalid': isSubmitting && !firstname}" />
                  <small class="p-invalid" v-if="isSubmitting && !firstname">Укажите имя клиента!</small>  
                </div>  
              </div>
              <div class="col-12 col-md-3">
                <div class="form-group">
                  <label for="middle-initial">Отчество</label>
                  <InputText :disabled="dataProcess==true" class="p-inputtext-sm form-control" id="middle-initial" 
                    required="false" maxlength="50" v-model.trim="middle_initial" />
                </div>  
              </div>
              <div class="col-12 col-md-2">
                <div class="form-group">
                  <label for="birthday">Дата рождения</label>
                  <InputMask :disabled="dataProcess==true" class="p-inputtext-sm form-control" id="birthday" 
                    required="true" mask="99.99.9999" v-model.trim="birthday"
                    :class="{'p-invalid': (isSubmitting && !birthday) || (isSubmitting && !isValidDate)}" />
                    <small class="p-invalid" v-if="isSubmitting && !birthday">Укажите дату рождения клиента!</small>  
                    <small class="p-invalid" v-if="isSubmitting && !isValidDate">Укажите корректную дату рождения клиента!</small>
                </div>  
              </div>  
            </div>  
            <div class="form-row">
              <div class="col-12 col-md-3">
                <div class="form-group">
                  <label for="inn">ИНН</label>
                  <InputText :disabled="dataProcess==true" class="p-inputtext-sm form-control" id="inn" 
                    required="false" maxlength="12" v-model.trim="inn" disabled/>
                </div>
              </div>
              <div class="col-12 col-md-9">
                <div class="form-group">
                  <label for="birth-place">Место рождения</label>
                  <InputText :disabled="dataProcess==true" class="p-inputtext-sm form-control" id="birth-place" 
                    required="true" maxlength="255" v-model.trim="birth_place" 
                    :class="{'p-invalid': isSubmitting && !birth_place}" />
                    <small class="p-invalid" v-if="isSubmitting && !birth_place">Укажите место рождения клиента!</small> 
                </div>  
              </div>
            </div>
            <div class="form-row">          
              <div class="col-12">
                <div class="form-group">
                  <label for="address">Адрес регистрации</label>
                  <InputText :disabled="dataProcess==true" class="p-inputtext-sm form-control" id="address" 
                    required="false" v-model.trim="address" /> 
                </div>  
              </div>
            </div>  
            <div class="form-row">
              <div class="col-12">
                <div class="form-group">
                  <label for="passport">Паспортные данные (серия, номер, кем и когда выдан, код подразделения)</label>
                  <InputText :disabled="dataProcess==true" class="p-inputtext-sm form-control" id="passport" 
                    required="true" v-model.trim="passport" 
                    :class="{'p-invalid': isSubmitting && !passport}" />
                  <small class="p-invalid" v-if="isSubmitting && !passport">Укажите паспортные данные клиента!</small> 
                </div>  
              </div>
            </div>  
            <div class="form-row">
              <div class="col-12 col-md-4">
                <div class="form-group">
                  <label for="phone1">Номер телефона №1</label>
                  <InputText :disabled="dataProcess==true" class="p-inputtext-sm form-control" id="phone1" 
                    required="false" maxlength="50" v-model.trim="phone1" disabled />
                </div>
              </div>
              <div class="col-12 col-md-4">
                <div class="form-group">
                  <label for="phone2">Номер телефона №2</label>
                  <InputText :disabled="dataProcess==true" class="p-inputtext-sm form-control" id="phone2" 
                    required="false" maxlength="50" v-model.trim="phone2" disabled />
                </div>
              </div>
              <div class="col-12 col-md-4">                
                <div class="form-group">
                  <label for="email">Адрес электронной почты</label>
                  <InputText :disabled="dataProcess==true" class="p-inputtext-sm form-control" id="email" 
                    required="false" maxlength="255" v-model.trim="email" disabled  />
                </div>
              </div>            
            </div>     

            <div class="form-row">          
              <div class="col-12">                
                <div class="form-group">
                  <label for="comment">Комментарий</label>
                  <Textarea :disabled="dataProcess==true" class="p-inputtext-sm form-control" id="comment" 
                    required="false" v-model="comment" /> 
                </div>  
              </div>
            </div>
            <div class="form-row">
              <div class="col-12">
                <div class="form-check form-check-inline">
              <RadioButton name="isMoscowRadioOptions" id="moscowRadioButton" :value="true" v-model="is_moscow"/>
              <label class="form-check-label" for="moscowRadioButton">&nbsp;<b>Москва и МО</b></label>
            </div>
            <div class="form-check form-check-inline">
              <RadioButton name="isMoscowRadioOptions" id="notMoscowRadioButton" :value="false" v-model="is_moscow"/>
              <label class="form-check-label" for="notMoscowRadioButton">&nbsp;<b>Регион</b></label>
            </div>
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
                            @click="current_scan = slotProps.data.img;display = true" />                                                      
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
           
      </div>        

    </div>

    <!-- просмотр скана -->
    <Dialog header="Скан экрана" :visible.sync="display" :modal="true" position="top">
	    <img :src="current_scan" style="max-width:1200px;max-height:750px;"/>
      <template #footer>
		    <Button label="Закрыть" icon="pi pi-times" class="p-button-text" @click="display = false"/>
	    </template>
    </Dialog>

    <!-- подтверждение браковки анкеты -->
    <Dialog :modal="true" header="Подтверждение" :visible.sync="defectDialogVisible" class="p-fluid">
      <i class="pi pi-exclamation-triangle p-mr-3" style="font-size: 2rem" />&nbsp;
      <span v-if="current_record">Вы уверены что хотите забраковать анкету клиента?</span>        
      <template #footer>
        <Button label="Не браковать" icon="pi pi-times" class="p-button-text" @click="defectDialogVisible = false"/>
        <Button label="Забраковать" icon="pi pi-check" class="p-button-text" @click="defectCustomer" />
      </template>
    </Dialog>

    <!-- подтверждение одобрения анкеты --->
    <Dialog v-if="list.length>0" :modal="true" header="Подтверждение" :visible.sync="confirmDialogVisible" class="p-fluid">
      <div class="confirmation-content">
        <i class="pi pi-exclamation-triangle p-mr-3" style="font-size: 2rem" />&nbsp;
        <span v-if="current_record">Вы уверены что хотите одобрить анкету клиента и сохранить данные?</span>
      </div>
      <template #footer>
        <Button label="Нет" icon="pi pi-times" class="p-button-text" @click="confirmDialogVisible = false"/>
        <Button label="Да" icon="pi pi-check" class="p-button-text" @click="confirmCustomer" />
      </template>
    </Dialog>
 
    </div>
</template>

<script lang="ts">
import Vue from 'vue';
const moment = require('moment');

import BlockUI from 'primevue/blockui';
import Button from 'primevue/button';
import Calendar from 'primevue/calendar';
import Card from 'primevue/card';
import Column from 'primevue/column';
import DataTable from 'primevue/datatable';
import DataView from 'primevue/dataview';
import Dialog from 'primevue/dialog';
import InputMask from 'primevue/inputmask';
import InputText from 'primevue/inputtext';
import Message from 'primevue/message';
import Panel from 'primevue/panel';
import RadioButton from 'primevue/radiobutton';
import Textarea from 'primevue/textarea';
import Toolbar from 'primevue/toolbar';

Vue.component('BlockUI', BlockUI);
Vue.component('Button', Button);
Vue.component('Calendar', Calendar);
Vue.component('Card', Card);
Vue.component('Column', Column);
Vue.component('DataTable', DataTable);
Vue.component('DataView', DataView);
Vue.component('Dialog', Dialog);
Vue.component('InputMask',InputMask);
Vue.component('InputText',InputText);
Vue.component('Message',Message);
Vue.component('Panel',Panel);
Vue.component('RadioButton',RadioButton);
Vue.component('Textarea', Textarea);
Vue.component('Toolbar', Toolbar);


export default Vue.extend({
  //данные страницы
    data() {
      return { 
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
        //список сканов экрана
        scans: [],
        //флаг отображения диалога просмотра скана
        display: false,
        //текущий скан для просмотра
        current_scan: "",
        //флаг отображения диалога подтверждения браковки анкеты
        defectDialogVisible: false,
        //флаг отображения диалога подтверждения одобрения анкеты
        confirmDialogVisible: false,
        //флаг валидности даты рождения
        isValidDate: true,
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
        //номера контактных телефонов
        phone1: "",
        phone2: "",
        //адрес электронной почты
        email: "",
        //флаг Москва / регионы
        is_moscow: null,
        //комментарий
        comment: "",
        
        filters: {},
        //------------------------------------
        //локализация календаря
        ru: {
          firstDayOfWeek: 1,
          dayNames: ["Воскресенье", "Понедельник", "Вторник", "Среда", "Четверг", "Пятница", "Суббота"],
          dayNamesShort: ["вс.", "пн.", "вт.", "ср.", "чт.", "пт.", "сб."],
          dayNamesMin: ["вс.", "пн.", "вт.", "ср.", "чт.", "пт.", "сб."],
          monthNames: [ "Январь","Февраль","Март","Апрель","Май","Июнь","Июль","Август","Сентябрь","Октябрь","Ноябрь","Декабрь" ],
          monthNamesShort: [ "янв", "фев", "мар", "апр", "май", "июн","июл", "авг", "сен", "окт", "ноя", "дек." ],
          today: 'сегодня',
          clear: 'очистить',
          dateFormat: 'dd.mm.yy',
          weekHeader: 'нед.'
        },
        //------------------------------------------------------
        filter_date: "",
        filter_list: [],
      }
    },    
    //метаданные страницы
    head() {        
        return {
            //заголовок браузера
            title: 'Неверифицированные анкеты | АРМ аудитора ИС Верификатор'
        }
    },    
    //загрузка статуса приложения
    async beforeCreate() {
       this.$store.commit('initializeStore');
       this.$store.commit('CHANGE_TITLE','Неверифицированные анкеты');
    },
    //--------------------------------------------------------
    //загрузка страницы    
    async mounted() {            
      this.errorMessageVisible = false;
      this.dataProcess = true;        
      try {
        const res = await this.$nuxt.$http.$post('/api/auditor/customers');
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
          (this.$refs.dt as DataTable.DataTable).selection = this.current_record;          
          //загружаем данные первой анкеты
          this.showCustomer();
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
        (this.$refs.dt as DataTable.DataTable).selection = this.current_record;                      
        this.showCustomer();    
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
        (this.$refs.dt as DataTable.DataTable).selection = this.current_record;                      
        this.showCustomer();                
      },
      //---------------------------------------------------
      //браковка анкеты клиента
      async defectCustomer() {              
          this.defectDialogVisible = false;
          //сохраняем данные
          let searchParams = new URLSearchParams();
          searchParams.set('id',(<any>this.current_record).id.toString());
          searchParams.set('userid',this.$store.state.userid.toString());
          searchParams.set('comment',this.comment);
          try {
            this.dataProcess = true;
            const res = await this.$nuxt.$http.$post('/api/auditor/defect-customer',searchParams);
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
            //удаляем забракованную запись из локального списка
            let index = 0;
            for (index=0; index<this.list.length; index++) {
              if ((<any>this.list[index]).id === (<any>this.current_record).id) {
                this.list.splice(index,1);
                if (index >= this.list.length) index--;
                break;
              }
            }
            if (this.list.length > 0) {              
              this.current_record = this.list[index];
              (this.$refs.dt as DataTable.DataTable).selection = this.current_record;
              //загружаем данные следующей в списке анкеты
              this.showCustomer();
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
      },
      
      //---------------------------------------------------
      //одобрение анкеты
      async confirmCustomer() {                
        this.confirmDialogVisible = false;
        this.isSubmitting = true;
        //проверяем валидность указанной даты рождения
        if (this.birthday.trim()) {
          if (!moment(this.birthday.trim(),"DD.MM.YYYY",true).isValid()) {
            this.isValidDate = false;
            return;
          } else {
            this.isValidDate = true;
          }
        }
        //проверяем заполнение всех обязательных полей 
        const required_fields = this.lastname.trim() && this.firstname.trim() && this.birthday.trim() 
          && this.birth_place.trim() && this.passport.trim();          
        if (required_fields)  {
          this.isSubmitting = false;
          //сохраняем данные
         // console.log(this.current_record)
          let searchParams = new URLSearchParams();
          searchParams.set('id',(<any>this.current_record).id.toString());
          searchParams.set('userid',this.$store.state.userid.toString());
          searchParams.set('lastname',this.lastname.trim());
          searchParams.set('firstname',this.firstname.trim());
          searchParams.set('middle_initial',this.middle_initial.trim());
          searchParams.set('birthday',this.birthday.trim());
          searchParams.set('birth_place',this.birth_place.trim());
          searchParams.set('address',this.address.trim());
          searchParams.set('passport',this.passport.trim());
          searchParams.set('comment',this.comment.trim());
          searchParams.set('is_moscow',this.is_moscow == null ? "null" : (<any>this.is_moscow).toString());
          try {
            this.dataProcess = true;
            const res = await this.$nuxt.$http.$post('/api/auditor/confirm-customer',searchParams);
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
            //удаляем одобренную запись из локального списка
            let index = 0;
            for (index=0; index<this.list.length; index++) {
              if ((<any>this.list[index]).id === (<any>this.current_record).id) {
                this.list.splice(index,1);
                if (index >= this.list.length) index--;
                break;
              }
            }
            if (this.list.length > 0) {              
              this.current_record = this.list[index];
              (this.$refs.dt as DataTable.DataTable).selection = this.current_record;
              //загружаем данные следующей в списке анкеты
              this.showCustomer();
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
        } 
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
          const res = await this.$nuxt.$http.$post('/api/auditor/get-customer',searchParams);
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
        this.comment = customer.comment == null ? "" : customer.comment;        
        this.is_moscow = customer.is_moscow;
        const phones = customer.phones;     
            
        this.phone1 = ""
        this.phone2 = ""
       if (phones != null) {
          this.phone1 = phones[0].phone;
          if (phones.length > 1) {
            this.phone2 = phones[1].phone;
          } 
        }       
        const emails = customer.emails;   
        if ((emails != null) && (emails.length > 0)) {
          this.email = emails[0].email;
        } else {
          this.email = ""
        }
        this.scans = customer.screen_scans;
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
      //--------------------------------------------------
      //фильтрация по дате
      async filterDate() {
        if (this.filter_date) {
          this.list = this.filter_list;
          this.list = this.list.filter(data => moment((<any>data).inserted).format("DD.MM.YYYY") ==  moment(this.filter_date).format("DD.MM.YYYY") )
        } else {
          this.list = this.filter_list;
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
