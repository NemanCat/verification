<template>
    <div>    
      <BlockUI :blocked="dataProcess" :fullScreen="true" class="block-ui">
        <img id="loading-image" src="~assets/img/ajax-loader.gif" alt="Загрузка страницы..." class="block-ui-image" v-if="dataProcess"/>
      </BlockUI>
      <Message v-if="errorMessageVisible == true" :severity="this.$store.state.fetchError.fetchErrorLevel==0 ? 'error' : 'warn'">{{this.$store.state.fetchError.fetchErrorMessage}}</Message>
      
        <Toolbar>
            <template #left>
                <form class="form-inline">
                    <label for="date-from" class="ml-1">Отчёт за период с</label>
                    <InputMask :disabled="dataProcess==true" class="form-control ml-3" id="date-from" 
                        required="true" mask="99.99.9999 99:99" v-model.trim="date_from" trim="Дата начала отчётного периода"
                        :class="{'p-invalid': (isSubmitting && !date_from) || (isSubmitting && !isValidDateFrom)}"
                        size="13" />
                    <small class="p-invalid" v-if="isSubmitting && !date_from">&nbsp;&nbsp;Укажите дату начала отчётного периода!</small>  
                    <small class="p-invalid" v-if="isSubmitting && !isValidDateFrom">&nbsp;&nbsp;Укажите корректную дату начала отчётного периода!</small>

                    <label for="date-to" class="ml-3">по</label>
                    <InputMask :disabled="dataProcess==true" class="form-control ml-3" id="date-to" 
                        required="true" mask="99.99.9999 99:99" v-model.trim="date_to" trim="Дата окончания отчётного периода"
                        :class="{'p-invalid': (isSubmitting && !date_to) || (isSubmitting && !isValidDateTo)}" 
                        size="13"/>
                    <small class="p-invalid" v-if="isSubmitting && !date_to">&nbsp;&nbsp;Укажите дату окончания отчётного периода!</small>  
                    <small class="p-invalid" v-if="isSubmitting && !isValidDateTo">&nbsp;&nbsp;Укажите корректную дату окончания отчётного периода!</small>

                    <Button type="button" icon="pi pi-table" label="Сформировать отчёт" 
                        title="Сформировать отчёт по обработанным паспортистом анкетам за указанный период" class="ml-4"
                        @click="makeReport" />                     
                </form>       
            </template>
        </Toolbar>    
        <DataTable :value="this.list" class="p-datatable-striped p-datatable-gridline" autoLayout>    
            <template #empty>
                За указанный период не обработано ни одной анкеты
            </template>        
            <Column field="name" header="Паспортист"></Column>           
            <Column field="verified_forms" header="Обработано паспортистом анкет"></Column>
            <Column field="defected_forms" header="Забраковано паспортистом анкет"></Column>
        </DataTable>  
    </div>
</template>

<script lang="ts">
    import Vue from 'vue'
    const moment = require('moment');

    import BlockUI from 'primevue/blockui';
    import Button from 'primevue/button';
    import Column from 'primevue/column';
    import DataTable from 'primevue/datatable';
    import InputMask from 'primevue/inputmask';
    import Message from 'primevue/message';
    import Toolbar from 'primevue/toolbar';
    

    Vue.component('BlockUI', BlockUI);
    Vue.component('Button', Button);
    Vue.component('Column', Column);
    Vue.component('DataTable', DataTable);
    Vue.component('InputMask',InputMask);
    Vue.component('Message',Message);
    Vue.component('Toolbar', Toolbar);

    export default Vue.extend({
        //данные страницы
        data() {
            return { 
                //флаг сохранения/загрузки данных
                dataProcess: false,
                //флаг отображения сообщения об ошибке
                errorMessageVisible: false,
                //флаг отправки формы
                isSubmitting: false,
                //даты начала и окончания отчётного периода
                date_from: null,
                date_to: null,
                //флаги валидности введённых дат
                isValidDateTo: true,
                isValidDateFrom: true,
                //данные отчёта
                list: [],
            }
        },
        //заголовки страницы 
        head() {
            return {
                //заголовок браузера
                title: 'Отчёт по паспортистам | АРМ аудитора ИС Верификатор',           
            }
        },
        async beforeCreate() {             
            this.$store.commit('initializeStore');
            this.$store.commit('CHANGE_TITLE','Отчёт по паспортистам');           
        }, 
        //-----------------------------------------------------------
        //загрузка страницы    
        async mounted() {
            //устанавливаем стартовые границы отчётного периода 
            //с первого числа текущего месяца по сегодняшний день
            this.date_to = moment().format("DD.MM.YYYY HH:hh");
            this.date_from = moment().startOf('month').format("DD.MM.YYYY HH:mm");            
            this.makeReport();
        },    
        //--------------------------------------------------------
        methods: {
            //формирование отчёта
            async makeReport() {
               // alert(moment((<any>this.date_from).trim(),"DD.MM.YYYY HH:mm").isValid())
                this.errorMessageVisible = false;
                this.isSubmitting = true;
                //проверяем валидность указанных границ отчётного периода
                if ((<any>this.date_from).trim()) {                    
                    if (!moment((<any>this.date_from).trim(),"DD.MM.YYYY HH:mm").isValid()) {
                        this.isValidDateFrom = false;
                        return;
                    } else {
                        this.isValidDateFrom = true;
                    }
                } else {
                    return
                }
                if ((<any>this.date_to).trim()) {
                    if (!moment((<any>this.date_to).trim(),"DD.MM.YYYY HH:mm").isValid()) {
                        this.isValidDateTo = false;
                        return;
                    } else {
                        this.isValidDateTo = true;
                    }
                } else {
                    return
                }
                //формируем отчёт
                const searchParams = new URLSearchParams();
                searchParams.set('userid',this.$store.state.userid.toString());
                searchParams.set('date_from',(<any>this.date_from).trim());
                searchParams.set('date_to',(<any>this.date_to).trim());
                this.isSubmitting = false;
                this.dataProcess = true;
                try {
                    const res = await this.$nuxt.$http.$post('/api/auditor/passport-report',searchParams);
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
                    //выводим отчёт
                    this.list = JSON.parse(res.Data);
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
        }        
    })
</script>    

<style>
.p-toolbar-group-left {
    width: 100%;
}
</style>