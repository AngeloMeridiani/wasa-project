<script>
import ChatsList from '@/components/ChatsList.vue';
import ChatWindow from '@/components/ChatWindow.vue';
export default {
	components : {
		ChatsList,
		ChatWindow
	},
	data: function() {
		return {
			errormsg: null,
			loading: false,
			availableEmojis: ['👍', '❤️', '💔', '😴', '😂', '😮', '😢', '🔥', '🤐']
		}
	},
	methods: {
		async refresh() {
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get("/");
				this.some_data = response.data;
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
	},
	mounted() {
		this.refresh()
	}
}
</script>

<template>
	<div class="row">
		
		<div class="col-md-4 border-end">
			<ChatsList />
		</div>
		
		<div class="col-md-8">
			<ChatWindow />
		</div>
		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
	</div>
</template>

<style>
</style>
