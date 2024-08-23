<template>
  <div class="container-fluid">
    <BlockUI :blocked="dataProcess" :fullScreen="true" class="block-ui">
        <img id="loading-image" src="~assets/img/ajax-loader.gif" alt="Загрузка страницы..." class="block-ui-image" v-if="dataProcess"/>
    </BlockUI>
    <Message v-if="errorMessageVisible" :severity="this.$store.state.fetchError.fetchErrorLevel==0 ? 'error' : 'warn'" >{{this.$store.state.fetchError.fetchErrorMessage}}</Message>  
    <Toolbar class="mb-2" >        
            <template #left>
              <Button label="Другой клиент" icon="pi pi-download" class="p-button-help mr-2" 
                title="Загрузить данные другого клиента для прозвона" 
                @click="getCustomer"/>  
              <Button label="Забраковать" icon="pi pi-times" class="p-button-danger" 
                title="Забраковать клиента как недоступного для контакта"    
                :disabled="!isSaveDataDisabled || isDefectButtonDisabled"           
                @click="defectDialogVisible = true" />                
            </template>        
            <template #right>        
              <Button label="Не готов к сотрудничеству" icon="pi pi-lock" class="p-button-danger mr-2" 
                title="Пометить клиента как не готового к сотрудничеству"
                :disabled="!isSaveDataDisabled || isSaveButtonsDisabled"
                @click="rejectCustomer" />
              <Button label="Готов к сотрудничеству" icon="pi pi-check" class="p-button-success" 
                title="Пометить клиента как готового к сотрудничеству" 
                :disabled="!isSaveDataDisabled || isSaveButtonsDisabled"
                @click="confirmCustomer" />          
            </template>
        </Toolbar> 
    <div v-show="this.list.length > 0">      
      <div class="row">
        <div class="col-5">
          <Panel>
            <template #header>
              Список анкет
            </template>
	          <DataTable :value="this.list" class="p-datatable-gridline p-datatable-sm"
              :paginator="true" :rows="10" :filters="filters"
              paginatorTemplate="CurrentPageReport FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink RowsPerPageDropdown"
              :rowsPerPageOptions="[10,20,50]"
              currentPageReportTemplate="Записи от {first} до {last} из {totalRecords}"
              @page="onPageSelect($event)"
              selectionMode="single" ref="dt" @row-select="onRowSelect">
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
              <Column field="phones_string" header="Номера телефонов" bodyClass="column-text-small" headerStyle="width: 20%;"
                filterField="phones_string" filterMatchMode="contains">
                <template  #body="slotProps">
                  <div style="text-align:justify;" v-html="slotProps.data.phones_string"></div>                        
                </template>
                <template #filter>
                  <InputText type="text" v-model="filters['phones_string']" class="p-column-filter p-inputtext-sm" size="12" placeholder="Телефон"/>
                </template> 
              </Column>
            </DataTable> 
          </Panel>          
        </div>
        <div class="col-7">
          <Panel>
            <template #header>
              Данные анкеты
            </template>
            <Toolbar v-if="list.length>0">                              
              <template #right>        
                <Button label="Сохранить" icon="pi pi-check" class="p-button-success p-mr-2" 
                  title="Сохранить внесённые изменения" :disabled="isSaveDataDisabled"
                  @click="saveDialogVisible=true;" />          
              </template>
            </Toolbar>
            <form v-if="list.length>0">
              <div class="form-row">
                <div class="col-12 col-md-4">
                  <div class="form-group">
                    <label for="lastname">Фамилия</label>
                    <InputText :disabled="dataProcess==true" class="p-inputtext-sm form-control" id="lastname" autofocus 
                      required="true" maxlength="50" v-model.trim="lastname"
                      :class="{'p-invalid': isSubmitting && !lastname}"
                      @change="enableSaveData" />
                    <small class="p-invalid" v-if="isSubmitting && !lastname">Укажите фамилию клиента!</small>
                  </div>  
                </div>
                <div class="col-12 col-md-3">
                  <div class="form-group">
                    <label for="firstname">Имя</label>
                    <InputText :disabled="dataProcess==true" class="p-inputtext-sm form-control" id="firstname" 
                      required="true" maxlength="50" v-model.trim="firstname"
                      :class="{'p-invalid': isSubmitting && !firstname}" 
                      @change="enableSaveData"/>
                    <small class="p-invalid" v-if="isSubmitting && !firstname">Укажите имя клиента!</small>  
                  </div>  
                </div>
                <div class="col-12 col-md-3">
                  <div class="form-group">
                    <label for="middle-initial">Отчество</label>
                    <InputText :disabled="dataProcess==true" class="p-inputtext-sm form-control" id="middle-initial" 
                      required="false" maxlength="50" v-model.trim="middle_initial" 
                      @change="enableSaveData"/>
                  </div>  
                </div>
                <div class="col-12 col-md-2">
                  <div class="form-group">
                    <label for="birthday">Дата рождения</label>
                    <InputMask :disabled="dataProcess==true" class="p-inputtext-sm form-control" id="birthday" 
                      required="true" mask="99.99.9999" v-model.trim="birthday"
                      :class="{'p-invalid': (isSubmitting && !birthday) || (isSubmitting && !isValidDate)}" 
                      @change="enableSaveData"/>
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
                      required="false" maxlength="12" v-model.trim="inn"
                      @change="enableSaveData"/>
                  </div>
                </div>
                <div class="col-12 col-md-9">
                  <div class="form-group">
                    <label for="birth-place">Место рождения</label>
                    <InputText :disabled="dataProcess==true" class="p-inputtext-sm form-control" id="birth-place" 
                      required="true" maxlength="255" v-model.trim="birth_place" 
                      :class="{'p-invalid': isSubmitting && !birth_place}" 
                      @change="enableSaveData"/>
                      <small class="p-invalid" v-if="isSubmitting && !birth_place">Укажите место рождения клиента!</small> 
                  </div>  
                </div>
              </div>
              <div class="form-row">          
                <div class="col-12">
                  <div class="form-group">
                    <label for="address">Адрес регистрации</label>
                    <InputText :disabled="dataProcess==true" class="p-inputtext-sm form-control" id="address" 
                      required="false" v-model.trim="address" 
                      @change="enableSaveData"/> 
                  </div>  
                </div>
              </div>  
              <div class="form-row">
                <div class="col-12">
                  <div class="form-group">
                    <label for="passport">Паспортные данные (серия, номер, кем и когда выдан, код подразделения)</label>
                    <InputText :disabled="dataProcess==true" class="p-inputtext-sm form-control" id="passport" 
                      required="true" v-model.trim="passport" 
                      :class="{'p-invalid': isSubmitting && !passport}" 
                      @change="enableSaveData"/>
                    <small class="p-invalid" v-if="isSubmitting && !passport">Укажите паспортные данные клиента!</small> 
                  </div>  
                </div>              
              </div>  
              <div class="form-row">          
                <div class="col-12">
                  <div class="form-group">
                    <label for="comment">Комментарий</label>
                    <Textarea :disabled="dataProcess==true" class="p-inputtext-sm form-control" id="comment" 
                      required="false" v-model.trim="comment" @change="enableSaveData"/> 
                  </div>  
                </div>
              </div>
              <div class="form-row">
                <div class="col-12">
                  <Panel>
                    <template #header>
                      Контакты клиента
                    </template>
                    <Button icon="pi pi-plus" class="p-button-rounded p-button-outlined mb-2" title="Добавить новый номер контактного телефона" @click="addNewPhone"></Button>
                    <DataTable ref="phones_dt" :value="this.phones" editMode="row" dataKey="id" :editingRows.sync="editingRows"
                      @row-edit-init="onRowEditInit" @row-edit-cancel="onRowEditCancel" @row-edit-save="onRowEditSave">            
                      <Column field="phone" header="Номер телефона">
                        <template #editor="slotProps">
                          <InputText v-model="slotProps.data[slotProps.column.field]" />
                        </template>
                      </Column>              
                      <Column field="result" header="Статус прозвона">
                        <template #body="slotProps">
                          {{call_results[slotProps.data.result].status}}
                        </template>
                        <template #editor="slotProps">
                          <Dropdown v-model="slotProps.data.result" :options="call_results" optionLabel="status" optionValue="result" style="height:100%;" />    
                        </template>    
                      </Column>
                      <Column field="call_datetime" header="Дата и время редактирования" />   
                      <Column  :rowEditor="true" headerStyle="width:7rem" bodyStyle="text-align:center" title="Редактировать данные"></Column>   
                      <Column headerStyle="width:3rem">
                        <template #body="slotProps">            
                          <Button icon="pi pi-trash" class="p-button-rounded p-button-text" title="Удалить номер телефона из списка контактных номеров" 
                            @click="current_phone=slotProps.data; confirmDeletePhoneDialogVisible=true;" />
                        </template>
                      </Column>             
                    </DataTable>
                    <hr class="mt-3" />
                    <Button icon="pi pi-plus" class="p-button-rounded p-button-outlined mb-2" title="Добавить новый адрес электронной почты" @click="addNewEmail"></Button>          
                    <DataTable ref="emails_dt" :value="this.emails" editMode="row"
                      dataKey="id" :editingRows.sync="editingEmailRows"
                      @row-edit-init="onEmailRowEditInit" @row-edit-cancel="onEmailRowEditCancel"
                      @row-edit-save="onEmailRowEditSave">            
                      <Column field="email" header="Адрес электронной почты">
                        <template #editor="slotProps">
                          <InputText v-model="slotProps.data[slotProps.column.field]" />
                      </template>
                      </Column>              
                      <Column  :rowEditor="true" headerStyle="width:7rem" bodyStyle="text-align:center" title="Редактировать данные"></Column>      
                      <Column headerStyle="width:3rem">
                        <template #body="slotProps">            
                          <Button icon="pi pi-trash" class="p-button-rounded p-button-text" 
                            title="Удалить адрес электронной почты из списка адресов" 
                            @click="current_email=slotProps.data; confirmDeleteEmailDialogVisible=true;" />
                        </template>
                      </Column>         
                    </DataTable>
                  </Panel>  
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
              <div class="form-row">
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
            <Toolbar v-if="list.length>0">                              
              <template #right>        
                <Button label="Сохранить" icon="pi pi-check" class="p-button-success p-mr-2" 
                  title="Сохранить внесённые изменения" :disabled="isSaveDataDisabled"
                  @click="saveDialogVisible=true;" />          
              </template>
            </Toolbar>       
          </Panel>        
        </div>        
      </div>  
    </div>  

    <!-- подтверждение удаления телефона --->
    <Dialog v-if="this.phones.length>0" :modal="true" header="Подтверждение" :visible.sync="confirmDeletePhoneDialogVisible" class="p-fluid">
      <div class="confirmation-content">
        <i class="pi pi-exclamation-triangle p-mr-3" style="font-size: 2rem" />&nbsp;
        <span v-if="current_phone != null">Вы уверены что хотите удалить номер <b>{{current_phone.phone}}</b> из списка номеров?</span>
      </div>
      <template #footer>
        <Button label="Нет" icon="pi pi-times" class="p-button-text" @click="confirmDeletePhoneDialogVisible = false"/>
        <Button label="Да" icon="pi pi-check" class="p-button-text" @click="deletePhone" />
      </template>
    </Dialog>

    <!-- подтверждение удаления адреса электронной почты --->
    <Dialog v-if="this.emails.length>0" :modal="true" header="Подтверждение" :visible.sync="confirmDeleteEmailDialogVisible" class="p-fluid">
      <div class="confirmation-content">
        <i class="pi pi-exclamation-triangle p-mr-3" style="font-size: 2rem" />&nbsp;
        <span v-if="current_email != null">Вы уверены что хотите удалить адрес электронной почты <b>{{current_email.email}}</b> из списка адресов?</span>
      </div>
      <template #footer>
        <Button label="Нет" icon="pi pi-times" class="p-button-text" @click="confirmDeleteEmailDialogVisible = false"/>
        <Button label="Да" icon="pi pi-check" class="p-button-text" @click="deleteEmail" />
      </template>
    </Dialog>

     <!-- подтверждение перехода на другую запись при несохранённых изменениях --->
    <Dialog v-if="this.list.length>0" :modal="true" header="Подтверждение" :visible.sync="confirmSelectionDialogVisible" class="p-fluid">
      <div class="confirmation-content">
        <i class="pi pi-exclamation-triangle p-mr-3" style="font-size: 2rem" />&nbsp;
        <span>ВНИМАНИЕ! В анкете имеются несохранённые изменения! Вы уверены что хотите загрузить другую анкету без сохранения изменений?</span>
      </div>
      <template #footer>
        <Button label="Нет" icon="pi pi-times" class="p-button-text" @click="confirmSelectionDialogVisible = false"/>
        <Button label="Да" icon="pi pi-check" class="p-button-text" @click="selectRow" />
      </template>
    </Dialog>

    <!-- подтверждение перехода на другую страницу пейджера при несохранённых изменениях --->
    <Dialog v-if="this.list.length>0" :modal="true" header="Подтверждение" :visible.sync="confirmPagerDialogVisible" class="p-fluid">
      <div class="confirmation-content">
        <i class="pi pi-exclamation-triangle p-mr-3" style="font-size: 2rem" />&nbsp;
        <span>ВНИМАНИЕ! В анкете имеются несохранённые изменения! Вы уверены что хотите загрузить другую анкету без сохранения изменений?</span>
      </div>
      <template #footer>
        <Button label="Нет" icon="pi pi-times" class="p-button-text" @click="confirmPagerDialogVisible = false"/>
        <Button label="Да" icon="pi pi-check" class="p-button-text" @click="selectPage" />
      </template>
    </Dialog>

    <!-- подтверждение сохранения анкеты -->
     <Dialog :modal="true" header="Подтверждение" :visible.sync="saveDialogVisible" class="p-fluid">
      <div class="confirmation-content">
        <i class="pi pi-exclamation-triangle p-mr-3" style="font-size: 2rem" />&nbsp;
        Вы уверены что все данные анкеты введены правильно?
      </div>
      <template #footer>
        <Button label="Нет" icon="pi pi-times" class="p-button-text" @click="saveDialogVisible = false"/>
        <Button label="Да" icon="pi pi-check" class="p-button-text" @click="saveData" />
      </template>
    </Dialog> 

     <!-- подтверждение браковки анкеты -->
      <Dialog :modal="true" header="Подтверждение" :visible.sync="defectDialogVisible" class="p-fluid">
        <div class="confirmation-content">
          <i class="pi pi-exclamation-triangle p-mr-3" style="font-size: 2rem" />&nbsp;
            Вы уверены что хотите забраковать данного клиента как недоступного для контакта?
        </div>
        <template #footer>
          <Button label="Нет" icon="pi pi-times" class="p-button-text" @click="defectDialogVisible = false"/>
          <Button label="Да" icon="pi pi-check" class="p-button-text" @click="defectCustomer" />
        </template>
      </Dialog>

  </div>
</template>

<script lang="ts">
import Vue from 'vue'
import  DataTable from 'primevue/datatable/DataTable';
const moment = require('moment');

export default Vue.extend({
  //данные страницы
    data() {
      return { 
        filters: {},
        //результаты прозвона
        call_results: [
          {"result": 0, "status": "не проверялся"},          
          {"result": 1, "status": "дозвонились"},  
          {"result": 2, "status": "номер предоставлен клиентом"},        
          {"result": 3, "status": "не берут или занято"},
          {"result": 4, "status": "номер не существует/не обслуживается"},
          {"result": 5, "status": "номер принадлежит другому человеку"},
        ],       
        //флаг сохранения/загрузки данных
        dataProcess: false,
        //флаг отображения сообщения об ошибке
        errorMessageVisible: false,
        //флаг процесса сохранения данных
        isSubmitting: false,        
        //список анкет
        list: [],
        //текущая запись
        current_record: [],
        //список сканов экрана оператора
        scans: [],
        //список сканов экрана паспортиста
        passport_scans: [],
        //список контактных телефонов
        phones: [],
        index: 0,
        originalRows: null,
        editingRows: [],
        current_phone: [],
        //список адресов электронной почты
        emails: [],
        emails_index: 0,
        originalEmailRows: null,
        editingEmailRows: [],
        current_email: [],        
        //текущий скан для просмотра
        current_scan: "",
        //флаг отображения диалога подтверждения браковки анкеты
        defectDialogVisible: false,
        //флаг отображения диалога подтверждения одобрения анкеты
        confirmDialogVisible: false,
        //флаг отображения диалога подтверждения сохранения внесённых изменений
        saveDialogVisible: false,
        //флаг валидности даты рождения
        isValidDate: true,
        //флаг отображения диалога подтверждения удаления номера телефона
        confirmDeletePhoneDialogVisible: false,
        //флаг отображения диалога подтверждения удаления адреса электронной почты
        confirmDeleteEmailDialogVisible: false,
        //флаг отображения подтверждения перехода на другую запись при несохранённых изменениях
        confirmSelectionDialogVisible: false,
        confirmPagerDialogVisible: false,
        //новый номер телефона
        new_phone_number: "",
        isNewPhoneSubmitting: false,
        //новый адрес электронной почты
        new_email: "",
        isNewEmailSubmitting: false,
        isSaveButtonsDisabled: true,
        isSaveDataDisabled: true,
        isDefectButtonDisabled: false,
        select_data_id: 0,
        first_id: 0,
        //--------------------------------------------------        
        //поля анкеты
        //id записи
        id: 0,
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
        //комментарий
        comment: "",
      }
    },    
    //метаданные страницы
    head() {        
        return {
            //заголовок браузера
            title: 'Анкеты клиентов | АРМ телефониста ИС Верификатор'
        }
    },    
    //загрузка статуса приложения
    async beforeCreate() {
       this.$store.commit('initializeStore');
       this.$store.commit('CHANGE_TITLE','Анкеты клиентов');
    },
    created() {        
      (<any>this.originalRows) = {};
      (<any>this.originalEmailRows) = {};
      window.addEventListener('beforeunload', this.beforeWindowUnload)
    },

    //--------------------------------------------------------
    //загрузка страницы    
    async mounted() {            
      this.errorMessageVisible = false;
      this.dataProcess = true;      
      const searchParams = new URLSearchParams();
      searchParams.set('userid',this.$store.state.userid.toString());  
      try {
        const res = await this.$nuxt.$http.$post('/api/assembly/get-customers-list',searchParams);
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
          for (let customer of this.list) {
            if ((<any>customer).phones != null) {
              for (let phone of (<any>customer).phones ) {                 
                 (<any>customer).phones_string += phone.phone + "<br>";   
              }  
            }
          }
          this.isSaveButtonsDisabled = false;
          this.current_record = this.list[0];
          (this.$refs.dt as DataTable).selection = this.current_record;           
          //загружаем данные первой анкеты
          this.showCustomer();
        } else {
          this.isDefectButtonDisabled = true;
          this.isSaveButtonsDisabled = true;
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
    //---------------------------------------------------------------------
    methods: {
      async showScanWindow() {
        const newTab = window.open();
        (<any>newTab).document.body.innerHTML = '<img src="' + this.current_scan + '" style="max-width:1200px;max-height:750px;">';
      },
      beforeWindowUnload(event: any) {
        if (!this.isSaveDataDisabled) {
          // Cancel the event
          event.preventDefault()
          // Chrome requires returnValue to be set
          event.returnValue = 'ВНИМАНИЕ! В анкете имеются несохранённые изменения!'
        }   
      },
      //выделение строки таблицы
      async onRowSelect(event: any) {
        this.select_data_id = event.data.id;
        if (!this.isSaveDataDisabled) {
          this.confirmSelectionDialogVisible = true;
        } else {
          this.selectRow();
        }
      },     
      async selectRow() {
        this.confirmSelectionDialogVisible = false;
        this.isSaveDataDisabled = true;
        const id: number = this.select_data_id;
        for (let i=0; i<this.list.length;i++) {
          if (this.list[i]['id'] === id) {
            this.current_record = this.list[i];
            break;
          }
        }                
        (this.$refs.dt as DataTable).selection = this.current_record;                      
        this.showCustomer();                
      }, 
      //выбор страницы в пейджере
      async onPageSelect(event: any) {   
        this.first_id = event.first_id;
        if (!this.isSaveDataDisabled) {
          this.confirmSelectionDialogVisible = true;
        } else {
          this.selectPage();
        }   
      },
      async selectPage() {
        this.confirmPagerDialogVisible = false;
        this.isSaveDataDisabled = true;
        //выбираем первую запись на странице        
        this.current_record = this.list[this.first_id];
        (this.$refs.dt as DataTable).selection = this.current_record;                      
        this.showCustomer(); 
      },
      //доступность кнопки сохранения изменений анкеты
      async enableSaveData() {        
        this.isSaveDataDisabled = false;
      },
      //------------------------------------------------------------------
      //заполнение формы данными текущей выделенной записи
      async showCustomer() {
        //выбираем данные анкеты
        let searchParams = new URLSearchParams();
        
        searchParams.set('id',(<any>this.current_record).id.toString());
        this.errorMessageVisible = false;
        try {
          this.dataProcess = true;
          const res = await this.$nuxt.$http.$post('/api/assembly/get-customer-data',searchParams);
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
        this.id = customer.id;        
        this.lastname = customer.lastname;
        this.firstname = customer.firstname;
        this.middle_initial = customer.middle_initial;
        this.birthday = moment(customer.birthday).format("DD.MM.YYYY");
        this.inn = customer.inn;
        this.birth_place = customer.birth_place;
        this.address = customer.address;
        this.passport = customer.passport_data;        
        this.scans = customer.screen_scans == null ? [] : customer.screen_scans;
        this.passport_scans = customer.passport_scans == null ? [] : customer.passport_scans;
        this.comment = customer.comment;
         //заполняем список контактных телефонов клиента
          this.isDefectButtonDisabled = true;
          this.isSaveButtonsDisabled = true;
          this.index = 0;          
          this.phones = [];
        
          for (let phone of customer.phones) {
            (<any>this.phones).push({"id": this.index, "phone": phone.phone, "result": phone.result == null ? 0 : phone.result, 
              "call_datetime": phone.call_datetime == null ? "" : phone.call_datetime });
              if (phone.result != null) {
                if (phone.result < 3) {
                   this.isDefectButtonDisabled = true;
                }
                if ((phone.result == 1) || (phone.result == 2)) {
                  this.isSaveButtonsDisabled = false;
                }
              }
            this.index++;              
          }
          
          //заполняем список адресов электронной почты клиента          
          this.emails_index = 0;
          this.emails = [];
          if (customer.emails != null) {
            for (let email of customer.emails) {
             (<any>this.emails).push({"id": this.emails_index, "email": email.email });
              this.emails_index++;              
            }
          }   
        
        this.new_phone_number = "";
        this.new_email = "";
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
      //-------------------------------------------------------------------------------
      //обработка редактирования данных телефона и прозвона
        async onRowEditInit(event: any) {
          (<any>this.originalRows)[event.index] = {...(<any>this.phones)[event.index]};
        },
        async onRowEditCancel(event: any) {     
          if ((<any>this.phones)[event.index].new) {
            (<any>this.phones).splice(event.index,1);
          } else {
            Vue.set((<any>this.phones), event.index, (<any>this.originalRows)[event.index]);
          }             
        },
        async onRowEditSave(event: any) {          
          (<any>this.phones)[event.index].new = false;
          //устанавливаем дату и время прозвона
          (<any>this.phones)[event.index].call_datetime = moment().format("DD.MM.YYYY HH:mm");
          this.isSaveDataDisabled = false;
          //разблокируем кнопки сохарнения данных
          let has_confirmed = await this.hasConfirmedPhone();
          if (has_confirmed) this.isSaveButtonsDisabled = false; else this.isSaveButtonsDisabled = true;
          //разблокируем кнопку браковки анкеты
          let enable_defect = await this.allPhonesDefected();
          if (enable_defect) this.isDefectButtonDisabled = false; else this.isDefectButtonDisabled = true;
        },
        //добавление нового контактного номера телефона
        async addNewPhone() {          
          (<any>this.phones).push({"id": this.index, "phone": "", "result": 0, "call_datetime": "", "new": true});
          (<any>this.originalRows)[this.index] = {...(<any>this.phones)[this.index]};
          (<any>(this.$refs.phones_dt as DataTable).editingRows).push((<any>this.phones)[this.index]);
          this.index++;
          this.isSaveDataDisabled = false;
        },
        //----------------------------------------------------------------------------------------
        //обработка редактирования данных электронной почты
        async onEmailRowEditInit(event: any) {
          (<any>this.originalEmailRows)[event.index] = {...(<any>this.emails)[event.index]};
        },        
        async onEmailRowEditCancel(event: any) {      
          if ((<any>this.emails)[event.index].new) {
            (<any>this.emails).splice(event.index,1);
          } else {
            Vue.set((<any>this.emails), event.index, (<any>this.originalEmailRows)[event.index]);
          }              
        },     
        async onEmailRowEditSave(event: any) {  
         this.isSaveDataDisabled = false;   
         (<any>this.emails)[event.index].new = false;
        },   
        //добавление нового адреса электронной почты
        async addNewEmail() {          
          (<any>this.emails).push({"id": this.emails_index, "email": "", "new": true});
          (<any>this.originalEmailRows)[this.emails_index] = {...(<any>this.emails)[this.emails_index]};
          (<any>(this.$refs.emails_dt as DataTable).editingRows).push((<any>this.emails)[this.emails_index]);
          this.emails_index++;
          this.isSaveDataDisabled = false;
        },
      //удаление номера телефона из списка
      async deletePhone() {
        this.confirmDeletePhoneDialogVisible = false;
        for (let i=0; i<this.phones.length; i++) {
          if ((<any>this.phones[i]).id === (<any>this.current_phone).id) {
            this.phones.splice(i--,1);
            break;
          }
        }
        let has_confirmed = await this.hasConfirmedPhone();
        if (!has_confirmed) this.isSaveButtonsDisabled = true;
        //разблокируем кнопку браковки анкеты
        let enable_defect = await this.allPhonesDefected();
        if (enable_defect) this.isDefectButtonDisabled = false; else this.isDefectButtonDisabled = true;
        this.isSaveDataDisabled = false;
      },  
      //удаление адреса электронной почты из списка
      async deleteEmail() {
        this.confirmDeleteEmailDialogVisible = false;
        for (let i=0; i<this.emails.length; i++) {
          if ((<any>this.emails[i]).id === (<any>this.current_email).id) {
            this.emails.splice(i--,1);
            break;
          }
        }
        this.isSaveDataDisabled = false;
      }, 
      //проверка наличия хотя бы одного достоверного телефона
      async hasConfirmedPhone(): Promise<boolean> {
        let has_confirmed = false;
        for (let phone of this.phones) {
          if ((<any>phone).result == 1) {
            has_confirmed = true;
            break;
          }
        }
        return has_confirmed;
      },
      //проверка все ли телефоны помечены как бракованные
      async allPhonesDefected(): Promise<boolean> {
        let all_defected = true;
        for (let phone of this.phones) {
          if ((<any>phone).result < 3) {
            all_defected = false;
            break;
          }
        }
        return all_defected;
      },  
      //--------------------------------------------------------------------------
      //сохранение данных анкеты
      async saveData() {
         this.saveDialogVisible = false;
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
          //формируем список контактных телефонов и адресов электронной почты
          let phones_list = [];
          for (let phone of (<any>this.phones)) {
            phones_list.push({"phone": phone.phone, "call_datetime": phone.call_datetime, "result": phone.result});
          } 
          let emails_list = [];
          for (let email of (<any>this.emails))   {
            emails_list.push({"email": email.email})
          }          
          //сохраняем данные          
          this.isSubmitting = false;
          this.$nuxt.$http.setHeader('Content-Type', 'application/json;charset=UTF-8')
          const request_body = {
            id: this.id.toString(),
            lastname: this.lastname.trim(),
            firstname: this.firstname.trim(),
            middle_initial: this.middle_initial.trim(),            
            birthday: this.birthday.trim(),
            inn: this.inn.trim(),
            birth_place: this.birth_place.trim(),
            address: this.address.trim(),
            passport: this.passport.trim(),
            phones: phones_list,
            emails: emails_list,
            comment: this.comment,
          }          
          this.errorMessageVisible = false;
          this.dataProcess = true;
          try {
            const res = await this.$nuxt.$http.$post('/api/assembly/save-customer',request_body);
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
            (<any>this.current_record).lastname = this.lastname;
            (<any>this.current_record).firstname = this.firstname;
            (<any>this.current_record).middle_initial = this.middle_initial;
            this.isSaveDataDisabled = true;
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
        return;
      },
      //-----------------------------------------------------------------------------
      //получение новой анкеты клиента на обработку
      async getCustomer() {
        this.errorMessageVisible = false;        
        this.dataProcess = true;    
        this.$nuxt.context.app.$http.setHeader('Content-Type', 'application/x-www-form-urlencoded');
        const searchParams = new URLSearchParams();
        searchParams.set('userid',this.$store.state.userid.toString());    
        try {
          const res = await this.$nuxt.$http.$post('/api/assembly/get-customer',searchParams);
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
          //загружаем данные анкеты нового клиента
          const customer = JSON.parse(res.Data)[0];     
          let phones_string = '';
          if (customer.phones != null) {
            for (let phone of customer.phones ) {
              phones_string += phone.phone + "<br>";
            }
          }
          
          (<any>this.list).push({
            id: customer.id,
            lastname: customer.lastname, 
            firstname: customer.firstname, 
            middle_initial: customer.middle_initial,
            phones: customer.phones,
            phones_string: phones_string});

          for (let i=0; i<this.list.length;i++) {
            if ((<any>this.list[i]).id == customer.id) {
              this.current_record = this.list[i];
              (this.$refs.dt as DataTable).selection = this.current_record; 
              break;
            }
          }  
          
          this.showCustomer();
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
      //------------------------------------------------------------------------------
      //браковка анкеты клиента
      async defectCustomer() {
        this.defectDialogVisible = false;
          //сохраняем данные
          let searchParams = new URLSearchParams();
          searchParams.set('id',this.id.toString());
          try {
            this.dataProcess = true;
            const res = await this.$nuxt.$http.$post('/api/assembly/defect-customer',searchParams);
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
            let rec_index = 0;
            for (rec_index=0; rec_index<this.list.length; rec_index++) {
              if ((<any>this.list[rec_index]).id === (<any>this.current_record).id) {
                this.list.splice(rec_index,1);
                if (rec_index >= this.list.length) rec_index--;
                break;
              }
            }
            if (this.list.length > 0) {                          
              this.current_record = this.list[rec_index];
              (this.$refs.dt as DataTable).selection = this.current_record;
              //загружаем данные следующей в списке анкеты
              this.showCustomer();
            }  else {
              this.isDefectButtonDisabled = true;
              this.isSaveButtonsDisabled = true;
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
      //------------------------------------------------------------------------------
      //клиент готов к дальнейшему сотрудничеству
      async confirmCustomer() {
        this.confirmDialogVisible = false;
          //сохраняем данные
          let searchParams = new URLSearchParams();
          searchParams.set('id',this.id.toString());
          try {
            this.dataProcess = true;
            const res = await this.$nuxt.$http.$post('/api/assembly/confirm-customer',searchParams);
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
            //удаляем  запись из локального списка
            let rec_index = 0;
            for (rec_index=0; rec_index<this.list.length; rec_index++) {
              if ((<any>this.list[rec_index]).id === (<any>this.current_record).id) {
                this.list.splice(rec_index,1);
                if (rec_index >= this.list.length) rec_index--;
                break;
              }
            }
            if (this.list.length > 0) {                          
              this.current_record = this.list[rec_index];
              (this.$refs.dt as DataTable).selection = this.current_record;
              //загружаем данные следующей в списке анкеты
              this.showCustomer();
            } else {
              this.isDefectButtonDisabled = true;
              this.isSaveButtonsDisabled = true;
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
      //------------------------------------------------------------------------------
      //клиент не готов к дальнейшему сотрудничеству
      async rejectCustomer() {
        this.confirmDialogVisible = false;
          //сохраняем данные
          let searchParams = new URLSearchParams();
          searchParams.set('id',this.id.toString());
          try {
            this.dataProcess = true;
            const res = await this.$nuxt.$http.$post('/api/assembly/reject-customer',searchParams);
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
            //удаляем  запись из локального списка
            let rec_index = 0;
            for (rec_index=0; rec_index<this.list.length; rec_index++) {
              if ((<any>this.list[rec_index]).id === (<any>this.current_record).id) {
                this.list.splice(rec_index,1);
                if (rec_index >= this.list.length) rec_index--;
                break;
              }
            }
            if (this.list.length > 0) {                          
              this.current_record = this.list[rec_index];
              (this.$refs.dt as DataTable).selection = this.current_record;
              //загружаем данные следующей в списке анкеты
              this.showCustomer();
            } else {
              this.isDefectButtonDisabled = true;
              this.isSaveButtonsDisabled = true;
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