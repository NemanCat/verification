<template>
  <div>    
      <BlockUI :blocked="dataProcess" :fullScreen="true" class="block-ui">
        <img id="loading-image" src="~assets/img/ajax-loader.gif" alt="Загрузка страницы..." class="block-ui-image" v-if="dataProcess"/>
      </BlockUI>
      <Message v-if="errorMessageVisible == true" :severity="this.$store.state.fetchError.fetchErrorLevel==0 ? 'error' : 'warn'">{{this.$store.state.fetchError.fetchErrorMessage}}</Message>
      <Toolbar>        
        <template #left>
          <Button label="Забраковать" icon="pi pi-times" class="p-button-danger" 
            title="Забраковать анкету и загрузить данные другого клиента из базы данных"
            :disabled="defectButtonDisabled"
            @click="defectDialogVisible = true" />
        </template>        
        <template #right>        
          <Button label="Сохранить и загрузить новую" icon="pi pi-check" class="p-button-success p-mr-2" 
             title="Сохранить введённые данные и загрузить данные другого клиента из базы данных" 
             @click="clone=false; saveDialogVisible=true;" />
          <!--   
          <Button label="Сохранить и клонировать" icon="pi pi-clone" class="p-button-help" 
            title="Сохранить введённые данные и создать новую анкету с исходными данными этого клиента"
            @click="clone=true; saveDialogVisible=true;" />
          -->  
        </template>
      </Toolbar>
      <!-- форма данных анкеты -->
      <form>
        <div class="form-row">
          <div class="col-12 col-md-3">
            <div class="form-group">
            <label for="lastname">Фамилия</label>
            <InputText :disabled="true" class="p-inputtext-sm form-control" id="lastname" autofocus 
              required="true" maxlength="50" v-model.trim="lastname"
              :class="{'p-invalid': isSubmitting && !lastname}" />
              <small class="p-invalid" v-if="isSubmitting && !lastname">Укажите фамилию клиента!</small>
            </div>  
          </div>
          <div class="col-12 col-md-3">
            <div class="form-group">
              <label for="firstname">Имя</label>
              <InputText :disabled="true" class="p-inputtext-sm form-control" id="firstname" 
                required="true" maxlength="50" v-model.trim="firstname"
                :class="{'p-invalid': isSubmitting && !firstname}" />
              <small class="p-invalid" v-if="isSubmitting && !firstname">Укажите имя клиента!</small>  
            </div>  
          </div>
          <div class="col-12 col-md-3">
            <div class="form-group">
              <label for="middle-initial">Отчество</label>
              <InputText :disabled="true" class="p-inputtext-sm form-control" id="middle-initial" 
                required="false" maxlength="50" v-model.trim="middle_initial" />
            </div>  
          </div>
          <div class="col-12 col-md-1">
            <div class="form-group">
              <label for="birthday">Дата рождения</label>
              <InputMask :disabled="dataProcess==true" class="p-inputtext-sm form-control" id="birthday" 
                required="true" mask="99.99.9999" v-model.trim="birthday"
                :class="{'p-invalid': (isSubmitting && !birthday) || (isSubmitting && !isValidDate)}" />
                <small class="p-invalid" v-if="isSubmitting && !birthday">Укажите дату рождения клиента!</small>  
                <small class="p-invalid" v-if="isSubmitting && !isValidDate">Укажите корректную дату рождения клиента!</small>
            </div>  
          </div>  
          <div class="col-12 col-md-2">
            <div class="form-group">
              <label for="inn">ИНН</label>
              <InputText :disabled="dataProcess==true" class="p-inputtext-sm form-control" id="inn" 
                required="false" maxlength="12" v-model.trim="inn" />
            </div>
          </div>
        </div>
        <div class="form-row">
          <div class="col-12 col-md-6">
            <div class="form-group">
              <label for="birth-place">Место рождения</label>
              <InputText :disabled="dataProcess==true" class="p-inputtext-sm form-control" id="birth-place" 
                required="true" maxlength="255" v-model.trim="birth_place" 
                :class="{'p-invalid': isSubmitting && !birth_place}" />
                <small class="p-invalid" v-if="isSubmitting && !birth_place">Укажите место рождения клиента!</small> 
            </div>  
          </div>
           <div class="col-12 col-md-6">
            <div class="form-group">
              <label for="address">Адрес регистрации</label>
              <InputText :disabled="dataProcess==true" class="p-inputtext-sm form-control" id="address" 
                required="false" v-model.trim="address" disabled /> 
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
                required="false" maxlength="255" v-model.trim="email" disabled />
            </div>
          </div>
        </div>      
        <div class="form-row">
          <div class="col-12">
            <div class="form-group">
              <label for="comment">Комментарий</label>
              <InputText :disabled="dataProcess==true" class="p-inputtext-sm form-control" id="comment" 
                required="true" v-model.trim="comment" 
                :class="{'p-invalid': (isSubmitting || isDefecting) && !comment}"/>
                <small class="p-invalid" v-if="(isSubmitting || isDefecting) && !comment">Укажите комментарий!</small> 
            </div>  
          </div>
        </div>
        <div class="form-row">
          <div class="col-12">
            <div class="form-check form-check-inline">
              <RadioButton name="isMoscowRadioOptions" id="moscowRadioButton" value="true" v-model="is_moscow"/>
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
                <Message v-if="typeErrorFlag==true" severity="error" >Недопустимый тип файла скана! Разрешены файлы с расширением jpg, jpeg, png. Укажите корректный файл скана!</Message>
                <Message v-if="sizeErrorFlag==true" severity="error" >Файл скана {{sizeErorFileName}} слишком большой! Размер файла скана не может превышать {{maxFileSize}} Мб.</Message>
                <div class="p-grid p-nogutter">
                  <div class="p-col-6" style="text-align: left">
                    Сканы экрана БД "Кронос"
                  </div>
                  <div class="p-col-6" style="text-align: right">
                    <Button label="Добавить сканы" icon="pi pi-plus" iconPos="left" title="Добавить новые сканы экрана в список" 
                      onclick="$('#file-input').trigger('click');"/>
                    <input id="file-input" type="file" multiple accept="image/png, image/jpeg" name="name" style="display: none;" @change="loadNewScan" />
                  </div>
                </div>
              </template>  
              <template #grid="slotProps">
                <div class="p-col-12 p-md-2">
                  <Card class="h-100">  
                    <template #header>                          
                        <Button icon="pi pi-search-plus" title="Просмотр скана" 
                          class="p-button-rounded float-left mb-2 ml-2" 
                          @click="current_scan = slotProps.data.img;display = true" />                            
                        <Button icon="pi pi-minus" title="Удалить скан из списка сканов" 
                          class="p-button-rounded p-button-danger float-right mb-2 mr-2"
                          @click="current_scan = slotProps.data.filename; deleteDialogVisible = true" />                      
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

      <!-- просмотр скана -->
      <Dialog header="Скан экрана" :visible.sync="display" :modal="true" class="p-fluid">
	      <img :src="current_scan"/>
        <template #footer>
		      <Button label="Закрыть" icon="pi pi-times" class="p-button-text" @click="display = false"/>
	      </template>
      </Dialog>

      <!-- подтверждение удаления скана -->
      <Dialog v-if="scans.length>0" :modal="true" header="Подтверждение" :visible.sync="deleteDialogVisible" class="p-fluid">
        <div class="confirmation-content">
          <i class="pi pi-exclamation-triangle p-mr-3" style="font-size: 2rem" />&nbsp;
          <span v-if="current_scan">Вы уверены что хотите удалить этот скан экрана?</span>
        </div>
        <template #footer>
          <Button label="Нет" icon="pi pi-times" class="p-button-text" @click="deleteDialogVisible = false"/>
          <Button label="Да" icon="pi pi-check" class="p-button-text" @click="removeScan" />
        </template>
      </Dialog>

      <!-- подтверждение браковки анкеты -->
      <Dialog :modal="true" header="Подтверждение" :visible.sync="defectDialogVisible" class="p-fluid">
        <div class="confirmation-content">
          <i class="pi pi-exclamation-triangle p-mr-3" style="font-size: 2rem" />&nbsp;
            Вы уверены что хотите забраковать эту анкету?
        </div>
        <template #footer>
          <Button label="Нет" icon="pi pi-times" class="p-button-text" @click="defectDialogVisible = false"/>
          <Button label="Да" icon="pi pi-check" class="p-button-text" @click="defectCustomer" />
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
     
    </div>  
      
</template>

<script lang="ts">
import Vue from 'vue'
const moment = require('moment');

const file2Base64 = (file:File):Promise<string> => {
    return new Promise<string> ((resolve,reject)=> {
         const reader = new FileReader();
         reader.readAsDataURL(file);
         reader.onload = () => resolve((<any>reader.result).toString());
         reader.onerror = error => reject(error);
     })
    }

export default Vue.extend({
   //данные страницы
    data() {
      return { 
        //флаг сохранения/загрузки данных
        dataProcess: false,
        //список сканов экрана
        scans: [],
        //флаг неправильного типа графического файла
        typeErrorFlag: false,
        //флаг отображения диалога просмотра скана
        display: false,
        //текущий скан для просмотра
        current_scan: "",
        //флаг отображения диалога подтверждения удаления скана
        deleteDialogVisible: false,
        //флаг отображения сообщения об ошибке
        errorMessageVisible: false,
        //флаги процесса сохранения данных
        isSubmitting: false,
        isDefecting: false,
        //флаг корректности даты
        isValidDate: true,
        //флаг отображения диалога подтверждения браковки анкеты
        defectDialogVisible: false,
        //флаг отображения диалога подтверждения браковки анкеты
        saveDialogVisible: false,
        clone: false,
        defectButtonDisabled: false,
        //--------------------------------------------------
        //максимальный размер файла скана, Мб
        maxFileSize: 5,
        //флаг отображения сообщения о превышении размера файла скана
        sizeErrorFlag: false,
        //имя слишком большого файла скана
        sizeErorFileName: "",
        //--------------------------------------------------
        //данные клиента из БД импорта
        data: [],
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
        //комментарий
        comment: "",
        //флаг Москва / регионы
        is_moscow: null,
      }
    },    
   //заголовки страницы 
   head() {
        return {
            //заголовок браузера
            title: 'Ввод анкеты клиента | АРМ оператора ИС Верификатор',           
        }
    },
    async beforeCreate() {             
       this.$store.commit('initializeStore');
       this.$store.commit('CHANGE_TITLE','Ввод анкеты клиента');           
    },  
    
    //--------------------------------------------------------
    //загрузка страницы    
    async mounted() {
      //загружаем данные импортированной анкеты
      const searchParams = new URLSearchParams();
      searchParams.set('userid',this.$store.state.userid.toString());
      this.errorMessageVisible = false;
      this.dataProcess = true;        
      try {
          const res = await this.$nuxt.$http.$post('/api/operator/import/get-customer',searchParams);
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
          this.data = JSON.parse(res.Data)[0];     
          this.$nuxt.$store.commit('CHANGE_CURRENTIMPORTID',(<any>this.data).id);
          this.lastname = (<any>this.data).lastname;
          this.firstname = (<any>this.data).firstname;
          this.middle_initial = (<any>this.data).middle_initial;
          this.birthday = "";
          this.inn = "";
          this.birth_place = "";
          this.address = (<any>this.data).address == null ? "" : (<any>this.data).address;
          this.passport = "";
          this.phone1 = (<any>this.data).phone1 == null ? "" : (<any>this.data).phone1;
          this.phone2 = (<any>this.data).phone2 == null ? "" : (<any>this.data).phone2;
          this.email = (<any>this.data).email == null ? "" : (<any>this.data).email;
          this.scans = [];
          this.comment = "";
          
          this.defectButtonDisabled = false;
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
    //--------------------------------------------------------
    methods: {
      //загрузка скана
      async loadNewScan(event: any) {
        //допустимые форматы файлов
        let available_types: string[] = [
          "image/jpeg",
          "image/png",          
        ];
        this.sizeErrorFlag = false;
        this.typeErrorFlag = false;
        for (let file of event.target.files) {
          //проверяем допустимость типа файла
          if (!available_types.includes(file.type)) {
            this.typeErrorFlag = true;
            return;
          }
          //ограничиваем макисмальный размер файла
          if (file.size > this.maxFileSize * 1024 * 1024) {
            this.sizeErorFileName = file.name;
            this.sizeErrorFlag = true;
            return;
          }
          let img = await file2Base64(file);
          //добавляем в локальный список сканов        
          (<any>this.scans).push({filename: file.name, img: img});
        }
      },
      //удаление скана из списка
      removeScan() {        
        this.deleteDialogVisible = false;
        //ищем скан в списке сканов и удаляем его
        let index = this.scans.findIndex(x => (<any>x).filename == this.current_scan);
        if (index > -1) {
          this.scans.splice(index,1);
        }        
      },
      //сохранение данных
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
          && this.birth_place.trim() && this.passport.trim() && this.comment.trim();
          
        if (required_fields)  {
          //проверяем наличие хотя бы одного скана
          if (this.scans.length == 0) {
            this.$store.commit('CHANGE_FETCHERROR', {
                fetchErrorFlag: true,
                fetchErrorCode: 0,
                fetchErrorMessage: "Добавьте к анкете клиента хотя бы один скан экрана БД Кронос!",
                fetchErrorLevel: 0,
              });
              this.errorMessageVisible = true;
              return;
          }
          //проверяем указан ли территориальный тип анкеты
          if (this.is_moscow == null) {
             this.$store.commit('CHANGE_FETCHERROR', {
                fetchErrorFlag: true,
                fetchErrorCode: 0,
                fetchErrorMessage: "Укажите, проживает ли клиент в Москве/МО или в регионе!",
                fetchErrorLevel: 0,
              });
              this.errorMessageVisible = true;
              return;
          }
          //сохраняем данные          
          this.isSubmitting = false;
          this.$nuxt.$http.setHeader('Content-Type', 'application/json;charset=UTF-8')
          const request_body = {
            importid: this.$store.state.currentImportId.toString(),
            lastname: this.lastname.trim(),
            firstname: this.firstname.trim(),
            middle_initial: this.middle_initial.trim(),
            birthday: this.birthday.trim(),
            inn: this.inn.trim(),
            birth_place: this.birth_place.trim(),
            address: this.address.trim(),
            passport: this.passport.trim(),
            phone1: this.phone1.trim(),
            phone2: this.phone2.trim(),
            email: this.email.trim(),
            scans: this.scans,
            operatorid: this.$store.state.userid.toString(),
            comment: this.comment,
            is_moscow: (<any>this.is_moscow).toString()
          }          
          
          this.errorMessageVisible = false;
          this.dataProcess = true;
          try {
            const res = await this.$nuxt.$http.$post('/api/operator/save-customer',request_body);
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

          if (this.clone) {
            //клонируем новую анкету из исходных данных           
            this.lastname = (<any>this.data).lastname;
            this.firstname = (<any>this.data).firstname;
            this.middle_initial = (<any>this.data).middle_initial;
            this.birthday = "";
            this.inn = "";
            this.birth_place = "";
            this.address = (<any>this.data).address == null ? "" : (<any>this.data).address;
            this.passport = "";
            this.phone1 = (<any>this.data).phone1 == null ? "" : (<any>this.data).phone1;
            this.phone2 = (<any>this.data).phone2 == null ? "" : (<any>this.data).phone2;
            this.email = (<any>this.data).email == null ? "" : (<any>this.data).email;
            this.scans = [];
            this.comment = "";
            this.is_moscow = null;
            this.defectButtonDisabled = true;
           } else {
            //загружаем нового клиента из таблицы импорта        
            this.$nuxt.context.app.$http.setHeader('Content-Type', 'application/x-www-form-urlencoded');       
            const searchParams = new URLSearchParams();
            searchParams.set('userid',this.$store.state.userid.toString());            
            this.errorMessageVisible = false;
            this.dataProcess = true;        
            try {
              const res = await this.$nuxt.$http.$post('/api/operator/import/get-customer',searchParams);
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
              this.data = JSON.parse(res.Data)[0];     
              this.$nuxt.$store.commit('CHANGE_CURRENTIMPORTID',(<any>this.data).id);
              this.lastname = (<any>this.data).lastname;
              this.firstname = (<any>this.data).firstname;
              this.middle_initial = (<any>this.data).middle_initial;
              this.birthday = "";
              this.inn = "";
              this.birth_place = "";
              this.address = (<any>this.data).address == null ? "" : (<any>this.data).address;
              this.passport = "";
              this.phone1 = (<any>this.data).phone1 == null ? "" : (<any>this.data).phone1;
              this.phone2 = (<any>this.data).phone2 == null ? "" : (<any>this.data).phone2;
              this.email = (<any>this.data).email == null ? "" : (<any>this.data).email;
              this.scans = [];
              this.comment = "";
              this.is_moscow = null;
              this.defectButtonDisabled = false;
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
          }
        }
      },      
      //браковка анкеты клиента
      async defectCustomer() {
        this.defectDialogVisible = false;        
        this.errorMessageVisible = false;        
        this.isDefecting = true;
        if (this.comment) {
          this.isDefecting = false;
          try {
          let searchParams = new URLSearchParams();
          searchParams.set('id',this.$store.state.currentImportId.toString());
          searchParams.set('comment',this.comment);
          this.dataProcess = true;
          let res = await this.$nuxt.$http.$post('/api/operator/import/defect-customer',searchParams);
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
          //анкета забракована, выбираем новую анкету
          searchParams = new URLSearchParams();          
          searchParams.set('userid',this.$store.state.userid.toString());
          this.dataProcess = true;
          res = await this.$nuxt.$http.$post('/api/operator/import/get-customer',searchParams);
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
          this.data = JSON.parse(res.Data)[0];     
          this.$nuxt.$store.commit('CHANGE_CURRENTIMPORTID',(<any>this.data).id);
          this.lastname = (<any>this.data).lastname;
          this.firstname = (<any>this.data).firstname;
          this.middle_initial = (<any>this.data).middle_initial;
          this.birthday = "";
          this.inn = "";
          this.birth_place = "";
          this.address = (<any>this.data).address;
          this.passport = "";
          this.phone1 = (<any>this.data).phone1;
          this.phone2 = (<any>this.data).phone2;
          this.email = (<any>this.data).email;
          this.comment = "";
          this.is_moscow = null;
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
