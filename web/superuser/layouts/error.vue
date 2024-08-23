<template>
  <div>
    <h1 v-if="error.statusCode === 404">Страница не найдена</h1>
    <h1 v-else>Ошибка приложения</h1>
    <p v-if="error.statusCode === 404">
        Запрошенная вами страница отсутствует в приложении. Проверьте, корректно ли указан адрес страницы.
    </p>
    <p v-else>
        Пожалуйста, сообщите разработчику, что в приложении произошла ошибка. Код ошибки: {{error.statusCode}}
    </p>
    <NuxtLink to="/">Стартовая страница приложения</NuxtLink>
  </div>
</template>

<script>
  export default {
    props: ['error'],
  }
</script>