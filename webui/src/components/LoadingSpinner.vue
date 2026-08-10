<script>
export default {
	props: ["loading"]
}
</script>

<template>
	<div v-if="loading" class="spinner-container">
		<div class="spinner-border text-primary" role="status">
			<span class="visually-hidden">Loading...</span>
		</div>
	</div>
	<div v-else class="content-container">
		<slot />
	</div>
</template>

<style scoped>
.spinner-container {
	display: flex;
	justify-content: center;
	align-items: center;
	width: 100%;
	height: 100%;
	min-height: calc(100vh - 60px); /* Altezza della finestra meno la navbar */
}
.content-container {
	width: 100%;
	height: 100%;
}
</style>
