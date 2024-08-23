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
                        required="true" mask="99.99.9999" v-model.trim="date_from" trim="Дата начала отчётного периода"
                        :class="{'p-invalid': (isSubmitting && !date_from) || (isSubmitting && !isValidDateFrom)}"
                        size="10" />
                    <small class="p-invalid" v-if="isSubmitting && !date_from">&nbsp;&nbsp;Укажите дату начала отчётного периода!</small>  
                    <small class="p-invalid" v-if="isSubmitting && !isValidDateFrom">&nbsp;&nbsp;Укажите корректную дату начала отчётного периода!</small>

                    <label for="date-to" class="ml-3">по</label>
                    <InputMask :disabled="dataProcess==true" class="form-control ml-3" id="date-to" 
                        required="true" mask="99.99.9999" v-model.trim="date_to" trim="Дата окончания отчётного периода"
                        :class="{'p-invalid': (isSubmitting && !date_to) || (isSubmitting && !isValidDateTo)}" 
                        size="10"/>
                    <small class="p-invalid" v-if="isSubmitting && !date_to">&nbsp;&nbsp;Укажите дату окончания отчётного периода!</small>  
                    <small class="p-invalid" v-if="isSubmitting && !isValidDateTo">&nbsp;&nbsp;Укажите корректную дату окончания отчётного периода!</small>

                    <Button type="button" icon="pi pi-table" label="Сформировать отчёт" 
                        title="Сформировать отчёт по введённым анкетам за указанный период" class="ml-4"
                        @click="makeReport" />                     
                </form>       
            </template>
        </Toolbar>    
        <DataTable v-if="list.length>0" :value="this.list" class="p-datatable-striped p-datatable-gridline" autoLayout>
            <Column field="total_customers" header="Всего анкет"></Column>
            <Column field="verified_customers" header="Верифицировано аудитором"></Column>
            <Column field="defected_customers" header="Забраковано аудитором"></Column>
        </DataTable>  
    </div>
</template>

<script lang="ts">
    import Vue from 'vue'
    const moment = require('moment');

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
                date_from: "",
                date_to: "",
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
                title: 'Отчёт по анкетам за период | АРМ оператора ИС Верификатор',           
            }
        },
        async beforeCreate() {             
            this.$store.commit('initializeStore');
            this.$store.commit('CHANGE_TITLE','Отчёт по анкетам за период');           
        }, 
        //-----------------------------------------------------------
        //загрузка страницы    
        async mounted() {
            //устанавливаем стартовые границы отчётного периода 
            //с первого числа текущего месяца по сегодняшний день
            this.date_to = moment().format("DD.MM.YYYY");
            this.date_from = moment().startOf('month').format("DD.MM.YYYY");
            this.makeReport();
        },    
        //--------------------------------------------------------
        methods: {
            //формирование отчёта
            async makeReport() {
                this.errorMessageVisible = false;
                this.isSubmitting = true;
                //проверяем валидность указанных границ отчётного периода
                if (this.date_from.trim()) {
                    if (!moment(this.date_from.trim(),"DD.MM.YYYY",true).isValid()) {
                        this.isValidDateFrom = false;
                        return;
                    } else {
                        this.isValidDateFrom = true;
                    }
                } else {
                    return
                }
                if (this.date_to.trim()) {
                    if (!moment(this.date_to.trim(),"DD.MM.YYYY",true).isValid()) {
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
                searchParams.set('date_from',this.date_from.trim());
                searchParams.set('date_to',this.date_to.trim());
                this.isSubmitting = false;
                this.dataProcess = true;
                try {
                    const res = await this.$nuxt.$http.$post('/api/operator/report',searchParams);
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