<template>
  <div class="container-fluid">
    <BlockUI :blocked="dataProcess" :fullScreen="true" class="block-ui">
        <img id="loading-image" src="~assets/img/ajax-loader.gif" alt="Загрузка страницы..." class="block-ui-image" v-if="dataProcess"/>
    </BlockUI>
    <Message v-if="errorMessageVisible" :severity="this.$store.state.fetchError.fetchErrorLevel==0 ? 'error' : 'warn'" >{{this.$store.state.fetchError.fetchErrorMessage}}</Message>  
    <Toolbar class="mb-2" >        
            <template #left>
              <Button label="Получить анкету" icon="pi pi-download" class="p-button-help mr-2" 
                title="Добавить новую анкету клиента в список анкет" 
                @click="loadCustomer" />                                
               <Button label="Вернуть анкету" icon="pi pi-upload" class="p-button-danger mr-2" 
                title="Вернуть анкету клиента в общую базу анкет" 
                :disabled="isSaveButtonsDisabled"
                @click="returnDialogVisible = true" /> 
            </template>        
            <template #right>        
              <Button label="Забраковать анкету" icon="pi pi-times" class="p-button-danger mr-2" 
                title="Забраковать анкету клиента как непригодную для дальнейшей обработки"    
                :disabled="isSaveButtonsDisabled"           
                @click="defectDialogVisible = true" />
              <Button label="Передать телефонисту" icon="pi pi-check" class="p-button-success" 
                title="Закончить работу с анкетой и передать её для дальнейшей обработки телефонисту" 
                :disabled="isSaveButtonsDisabled"
                @click="confirmDialogVisible = true;" />          
            </template>
        </Toolbar>
    <div v-show="this.list.length > 0">
      
      <div class="row">
        <div class="col-5">
          <Panel>
            <template #header>
              Список анкет
            </template>
	          <DataTable :value="this.list" class="p-datatable-gridline p-datatable-sm" autoLayout 
              :paginator="true" :rows="10"
              paginatorTemplate="CurrentPageReport FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink RowsPerPageDropdown"
              :rowsPerPageOptions="[10,20,50]"
              currentPageReportTemplate="Записи от {first} до {last} из {totalRecords}"
              @page="onPageSelect($event)"
              selectionMode="single" ref="dt" @row-select="onRowSelect">
              <template #empty>
                Нет ни одной записи
              </template>          
              <Column field="lastname" header="Фамилия" bodyClass="column-text-small"></Column>
              <Column field="firstname" header="Имя" bodyClass="column-text-small"></Column>
              <Column field="middle_initial" header="Отчество" bodyClass="column-text-small"></Column>              
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
                <div class="col-6">
                  <fieldset class="border p-2">
                    <legend class="w-auto">Список контактных телефонов</legend>     
                    <div class="p-inputgroup mb-1">
                      <InputText v-model.trim="new_phone_number" placeholder="Новый номер телефона"
                        :class="{'p-invalid': isNewPhoneSubmitting && !new_phone_number}"
                         @keyup.enter="addNewPhone" />
                      <Button icon="pi pi-plus" title="Добавить новый номер контактного телефона в список номеров"  @click="addNewPhone"/>
                    </div>
                    <small class="p-invalid" v-if="isNewPhoneSubmitting && !new_phone_number">Укажите номер телефона!</small> 
                    <DataTable ref="phones_dt" :value="this.phones" editMode="row" dataKey="id" :editingRows.sync="editingRows"
                      @row-edit-init="onRowEditInit" @row-edit-cancel="onRowEditCancel" @row-edit-save="onRowEditSave" class="mb-2">            
                      <Column field="phone" header="Номер телефона">
                        <template #editor="slotProps">
                          <InputText v-model="slotProps.data[slotProps.column.field]" />
                        </template>
                      </Column>              
                      <Column  :rowEditor="true" headerStyle="width:7rem" bodyStyle="text-align:center" title="Редактировать данные"></Column>   
                      <Column headerStyle="width:3rem">
                        <template #body="slotProps">            
                          <Button icon="pi pi-trash" class="p-button-rounded p-button-text" title="Удалить номер телефона из списка контактных номеров" 
                            @click="current_phone=slotProps.data; confirmDeletePhoneDialogVisible=true;" />
                        </template>
                      </Column>             
                    </DataTable>
                  </fieldset>
                </div>  
                <div class="col-6">
                  <fieldset class="border p-2">
                    <legend class="w-auto">Список адресов электронной почты</legend>                    
                      <div class="p-inputgroup mb-1">
                        <InputText v-model.trim="new_email" placeholder="Новый адрес электронной почты"
                          :class="{'p-invalid': isNewEmailSubmitting && !new_email}"
                          @keyup.enter="addNewEmail" />
                        <Button icon="pi pi-plus" title="Добавить новый адрес электронной почты в список адресов"  @click="addNewEmail"/>
                      </div>
                      <small class="p-invalid" v-if="isNewEmailSubmitting && !new_email">Укажите адрес электронной почты!</small>
                      <DataTable ref="emails_dt" :value="this.emails" editMode="row"
                        dataKey="id" :editingRows.sync="editingEmailRows"
                        @row-edit-init="onEmailRowEditInit" @row-edit-cancel="onEmailRowEditCancel"  @row-edit-save="onEmailRowEditSave">            
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
                  </fieldset>
                </div>
              </div>                      
              <div class="form-row">
                <div class="col-12">
                  <DataView :value="passport_scans" layout="grid" :paginator="false">
                    <template #empty>Нет ни одного скана</template>
                    <template #header>     
                      <Message v-if="typeErrorFlag==true" severity="error" >Недопустимый тип файла скана! Разрешены файлы с расширением jpg, jpeg, png. Укажите корректный файл скана!</Message>
                      <Message v-if="sizeErrorFlag==true" severity="error" >Файл скана {{sizeErorFileName}} слишком большой! Размер файла скана не может превышать {{maxFileSize}} Мб.</Message>           
                      <div class="p-grid p-nogutter">
                        <div class="p-col-6" style="text-align: left">
                          Мои сканы
                        </div>  
                        <div class="p-col-6" style="text-align: right">
                          <Button label="Добавить сканы" icon="pi pi-plus" iconPos="left" title="Добавить новые сканы экрана в список" 
                            onclick="$('#file-input').trigger('click');"/>
                          <input id="file-input" type="file" multiple accept="image/png, image/jpeg" name="name" style="display: none;" @change="loadNewScan" />
                        </div>                
                      </div>
                    </template>  
                    <template #grid="slotProps">
                      <div class="p-col-12 p-md-4">
                        <Card class="h-100">  
                          <!--
                          <template #header>                          
                            <Button icon="pi pi-search-plus" title="Просмотр скана" 
                              class="p-button-rounded float-left mb-2 ml-2"                               
                              @click="current_scan = slotProps.data.img;display = true" />                                                      
                          </template>
                          -->
                          <template #header>                          
                            <Button icon="pi pi-search-plus" title="Просмотр скана" 
                              class="p-button-rounded float-left mb-2 ml-2"                               
                              @click="current_scan = slotProps.data.img; showScanWindow();" />   
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
                          <!--
                          <template #header>                          
                            <Button icon="pi pi-search-plus" title="Просмотр скана" 
                              class="p-button-rounded float-left mb-2 ml-2"                               
                              @click="current_scan = slotProps.data.img;display = true" />                                                      
                          </template>
                          -->
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

    <!-- просмотр скана -->
    <Dialog header="Скан экрана" :visible.sync="display" :modal="true" position="top">
	    <img :src="current_scan" style="max-width:1200px;max-height:750px;"/>
      <template #footer>
		    <Button label="Закрыть" icon="pi pi-times" class="p-button-text" @click="display = false"/>
	    </template>
    </Dialog>

    <!-- подтверждение возврата анкеты --->
    <Dialog v-if="list.length>0" :modal="true" header="Подтверждение" :visible.sync="returnDialogVisible" class="p-fluid">
      <div class="confirmation-content">
        <i class="pi pi-exclamation-triangle p-mr-3" style="font-size: 2rem" />&nbsp;
        <span v-if="current_record">
          <span v-if="!isSaveDataDisabled">ВНИМАНИЕ! В анкете имеются несохранённые изменения!&nbsp;</span>
          Вы уверены что хотите вернуть анкету клиента <b>{{this.current_record.lastname}}&nbsp;{{this.current_record.firstname}}</b> в общую базу?
        </span>
      </div>
      <template #footer>
        <Button label="Нет" icon="pi pi-times" class="p-button-text" @click="returnDialogVisible = false"/>
        <Button label="Да" icon="pi pi-check" class="p-button-text" @click="returnCustomer" />
      </template>
    </Dialog>

    <!-- подтверждение браковки анкеты --->
    <Dialog v-if="list.length>0" :modal="true" header="Подтверждение" :visible.sync="defectDialogVisible" class="p-fluid">
      <div class="confirmation-content">
        <i class="pi pi-exclamation-triangle p-mr-3" style="font-size: 2rem" />&nbsp;
        <span v-if="current_record">
            <span v-if="!isSaveDataDisabled">ВНИМАНИЕ! В анкете имеются несохранённые изменения!&nbsp;</span>
            Вы уверены что хотите забраковать анкету клиента <b>{{this.current_record.lastname}}&nbsp;{{this.current_record.firstname}}</b>?
        </span>
      </div>
      <template #footer>
        <Button label="Нет" icon="pi pi-times" class="p-button-text" @click="defectDialogVisible = false"/>
        <Button label="Да" icon="pi pi-check" class="p-button-text" @click="defectCustomer" />
      </template>
    </Dialog>

    <!-- подтверждение одобрения анкеты --->
    <Dialog v-if="list.length>0" :modal="true" header="Подтверждение" :visible.sync="confirmDialogVisible" class="p-fluid">
      <div class="confirmation-content">
        <i class="pi pi-exclamation-triangle p-mr-3" style="font-size: 2rem" />&nbsp;
        <span v-if="current_record">
          <span v-if="!isSaveDataDisabled">ВНИМАНИЕ! В анкете имеются несохранённые изменения!&nbsp;</span>
          Вы уверены что хотите передать телефонистам анкету клиента <b>{{this.current_record.lastname}}&nbsp;{{this.current_record.firstname}}</b>?
        </span>
      </div>
      <template #footer>
        <Button label="Нет" icon="pi pi-times" class="p-button-text" @click="confirmDialogVisible = false"/>
        <Button label="Да" icon="pi pi-check" class="p-button-text" @click="confirmCustomer" />
      </template>
    </Dialog>

    <!-- подтверждение сохранения внесённых изменений -->
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

  </div>
</template>

<script lang="ts">
import Vue from 'vue'
import  DataTable from 'primevue/datatable/DataTable';
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
        //флаг отображения диалога просмотра скана
        display: false,
        //текущий скан для просмотра
        current_scan: "",
        //флаг отображения диалога подтверждения возврата анкеты
        returnDialogVisible: false,
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
        select_data_id: 0,
        first_id: 0,
        //флаг неправильного типа графического файла
        typeErrorFlag: false,
        //максимальный размер файла скана, Мб
        maxFileSize: 5,
        //флаг отображения сообщения о превышении размера файла скана
        sizeErrorFlag: false,
        //имя слишком большого файла скана
        sizeErorFileName: "",
        //флаг отображения диалога подтверждения удаления скана
        deleteDialogVisible: false,
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
            title: 'Анкеты клиентов | АРМ паспортиста ИС Верификатор'
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
        const res = await this.$nuxt.$http.$post('/api/passport/get-customers-list',searchParams);
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
                     
          this.isSaveButtonsDisabled = false;
          this.current_record = this.list[0];
          (this.$refs.dt as DataTable).selection = this.current_record;           
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
    //---------------------------------------------------------------------
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
          (<any>this.passport_scans).push({filename: file.name, img: img});
          this.isSaveDataDisabled = false;
        }
      },
      //удаление скана из списка
      removeScan() {        
        this.deleteDialogVisible = false;
        //ищем скан в списке сканов и удаляем его
        let index = this.passport_scans.findIndex(x => (<any>x).filename == this.current_scan);
        if (index > -1) {
          this.passport_scans.splice(index,1);
          this.isSaveDataDisabled = false;
        }                
      },
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
        
      //доступность кнопки сохранения изменений анкеты
      async enableSaveData() {        
        this.isSaveDataDisabled = false;
      },
      //----------------------------------------------------------------------------
      //получение новой анкеты клиента
      async loadCustomer() {
        this.errorMessageVisible = false;
        this.dataProcess = true;      
        const searchParams = new URLSearchParams();
        searchParams.set('userid',this.$store.state.userid.toString());  
        try {
          const res = await this.$nuxt.$http.$post('/api/passport/get-customer',searchParams);
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
          const customer = JSON.parse(res.Data)[0];     
          (<any>this.list).push(customer);
          this.current_record = customer;
          (this.$refs.dt as DataTable).selection = this.current_record;                      
          this.isSaveButtonsDisabled = false;
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
      //браковка анкеты клиента
      async defectCustomer() {         
        this.isSubmitting = true;        
          this.defectDialogVisible = false;
          this.isSubmitting = false;
          //сохраняем данные
          let searchParams = new URLSearchParams();
          searchParams.set('id',this.id.toString());
          try {
            this.dataProcess = true;
            const res = await this.$nuxt.$http.$post('/api/passport/defect-customer',searchParams);
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
      //одобрение анкеты клиента
      async confirmCustomer() {         
        this.confirmDialogVisible = false;
        //проверяем наличие хотя бы одного скана           
        if (this.passport_scans.length == 0) {
          this.$store.commit('CHANGE_FETCHERROR', {
            fetchErrorFlag: true,
            fetchErrorCode: 0,
            fetchErrorMessage: "Добавьте к анкете клиента хотя бы один скан экрана с подтверждением номера телефона!",
            fetchErrorLevel: 0,
          });
          this.errorMessageVisible = true;
          return;
        }        
        //проверяем наличие хотя бы одного номера телефона           
        if (this.phones.length == 0) {
          this.$store.commit('CHANGE_FETCHERROR', {
            fetchErrorFlag: true,
            fetchErrorCode: 0,
            fetchErrorMessage: "Добавьте к анкете клиента хотя бы один номер телефона!",
            fetchErrorLevel: 0,
          });
          this.errorMessageVisible = true;
          return;
        }
        //сохраняем данные
        let searchParams = new URLSearchParams();
        searchParams.set('id',this.id.toString());
        try {
            this.dataProcess = true;
            const res = await this.$nuxt.$http.$post('/api/passport/confirm-customer',searchParams);
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
      //заполнение формы данными текущей выделенной записи
      async showCustomer() {
        //выбираем данные анкеты
        let searchParams = new URLSearchParams();
        searchParams.set('id',(<any>this.current_record).id.toString());
        this.errorMessageVisible = false;
        try {
          this.dataProcess = true;
          const res = await this.$nuxt.$http.$post('/api/passport/get-customer-data',searchParams);
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
        this.index = 0;          
        this.phones = [];
        for (let phone of customer.phones) {
          (<any>this.phones).push({"id": this.index, "phone": phone.phone});
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
      //возврат анкеты клиента в базу
      async returnCustomer() {         
          this.returnDialogVisible = false;
          //сохраняем данные
          let searchParams = new URLSearchParams();
          searchParams.set('id',(<any>this.current_record).id.toString());
          try {
            this.dataProcess = true;
            const res = await this.$nuxt.$http.$post('/api/passport/return-customer',searchParams);
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
            //удаляем запись из локального списка
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
              //загружаем данные первой анкеты
              this.showCustomer();
            } else {
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
      //сохранение данных анкеты
      async saveData() {
         this.isSubmitting = true;
         this.saveDialogVisible = false;
        //проверяем валидность указанной даты рождения
        if (this.birthday.trim()) {
          if (!moment(this.birthday.trim(),"DD.MM.YYYY",true).isValid()) {
            this.isValidDate = false;
            return false;
          } else {
            this.isValidDate = true;
          }
        }
        //проверяем заполнение всех обязательных полей 
        const required_fields = this.lastname.trim() && this.firstname.trim() && this.birthday.trim() 
          && this.birth_place.trim() && this.passport.trim();
        
        if (required_fields)  {      
          //проверяем наличие хотя бы одного скана           
          if (this.passport_scans.length == 0) {
            this.$store.commit('CHANGE_FETCHERROR', {
                fetchErrorFlag: true,
                fetchErrorCode: 0,
                fetchErrorMessage: "Добавьте к анкете клиента хотя бы один скан экрана с подтверждением номера телефона!",
                fetchErrorLevel: 0,
              });
              this.errorMessageVisible = true;
              return;
          }
          //формируем список контактных телефонов и адресов электронной почты
          let phones_list = [];
          for (let phone of (<any>this.phones)) {
            phones_list.push({"phone": phone.phone});
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
            userid: this.$store.state.userid.toString(),
            status: 3,
            comment: this.comment,
            scans: this.passport_scans
          }          
          this.errorMessageVisible = false;
          this.dataProcess = true;
          try {
            const res = await this.$nuxt.$http.$post('/api/passport/save-customer',request_body);
            this.dataProcess = false;
            if(!res.Success) {              
              this.$store.commit('CHANGE_FETCHERROR', {
                fetchErrorFlag: true,
                fetchErrorCode: 0,
                fetchErrorMessage: res.Message,
                fetchErrorLevel: res.Level,
              });
              this.errorMessageVisible = true;
              return
            }
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
            return
          }  
        }          
      },
      //-------------------------------------------------------------------------------
      //обработка редактирования данных телефона и прозвона
        async onRowEditInit(event: any) {
          (<any>this.originalRows)[event.index] = {...(<any>this.phones)[event.index]};
        },
        async onRowEditCancel(event: any) {   
            Vue.set((<any>this.phones), event.index, (<any>this.originalRows)[event.index]);
        },
        async onRowEditSave(event: any) {  
          this.isSaveDataDisabled = false;
        },
        //добавление нового контактного номера телефона
        async addNewPhone() {                   
          this.isNewPhoneSubmitting = true;
          if (this.new_phone_number) {
            this.isNewPhoneSubmitting = false;
            (<any>this.phones).push({"id": this.index, "phone": this.new_phone_number});
            this.isSaveDataDisabled = false;
            this.index++;
            this.new_phone_number = "";
          }
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
          this.isSaveDataDisabled = false;
        },  
        //----------------------------------------------------------------------------------------
        //обработка редактирования данных электронной почты
        async onEmailRowEditInit(event: any) {
          (<any>this.originalEmailRows)[event.index] = {...(<any>this.emails)[event.index]};          
        },        
        async onEmailRowEditCancel(event: any) {                
            Vue.set((<any>this.emails), event.index, (<any>this.originalEmailRows)[event.index]);
        },     
        async onEmailRowEditSave(event: any) {                                     
          this.isSaveDataDisabled = false;
        },   
        //добавление нового адреса электронной почты
        async addNewEmail() {      
          this.isNewEmailSubmitting = true;
          if (this.new_email) {
            this.isNewEmailSubmitting = false;
            (<any>this.emails).push({"id": this.emails_index, "email": this.new_email});
            this.isSaveDataDisabled = false;
            this.emails_index++;
            this.new_email = "";
          }
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