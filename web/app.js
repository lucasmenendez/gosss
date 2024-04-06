import { createApp } from 'https://unpkg.com/vue@3/dist/vue.esm-browser.prod.js'

const app = createApp({
    data() {
        return {
            config: {
                minShares: 0,
                maxShares: 0,
                minMin: 0,
                maxMin: 0,
            },
            message: "",
            sharesCount: 0,
            threshold: 0,
            shares: "",
            currentTab: "hide",
            hide_result: "",
            recovered_message: "",
        }
    },
    async created() {
        await this.setupWebAssembly();
    },
    template: `
        <div>
            <button class="button" @click="currentTab = 'hide'">Hide</button>
            <button class="button" @click="currentTab = 'recover'">Recover</button>
            <div v-show="currentTab === 'hide'">
                <h3>Hide a message</h3>
                <textarea class="textarea" v-model="message" placeholder="Enter your message" rows="6"></textarea>
                <div v-show="message !== ''">
                    <label class="label">Shares <small>({{ config.minShares }} - {{ config.maxShares }})</small></label>
                    <input v-model="sharesCount" type="number" class="input" :min="config.minShares" :max="config.maxShares" :step="config.minShares">
                </div>
                <div v-show="message !== ''">
                    <label class="label">Threshold <small>({{ config.minMin }} - {{ sharesCount - 1 }})</small></label>
                    <input v-model="threshold" type="number" class="input" :min="config.minMin" :max="sharesCount - 1">
                </div>
                <button class="button is-secondary" :class="{'is-disabled': message == ''}" @click="hideMessage">Hide</button>
                <textarea class="textarea" readonly v-model="hide_result" placeholder="Resulting secret parts" rows="6"></textarea>
            </div>
            <div v-show="currentTab === 'recover'">
                <h3>Recover a message</h3>
                <textarea class="textarea" v-model="shares" placeholder="Enter shares" rows="6"></textarea>
                <button class="button is-secondary" :class="{'is-disabled': message == ''}" @click="recoverMessage">Recover</button>

                <textarea class="textarea" readonly v-model="recovered_message" placeholder="Recovered message" rows="6"></textarea>
            </div>
        </div>
    `,
    methods: {
        async setupWebAssembly() {
            const go = new Go();
            const result = await WebAssembly.instantiateStreaming(fetch("gosss.wasm"), go.importObject);
            go.run(result.instance);
        },
        getConfig() {
            const rawResult = GoSSS.limits(this.message);
            const result = JSON.parse(rawResult);
            if (!result.error) {
                this.config = {
                    minShares: result.data[0],
                    maxShares: result.data[1],
                    minMin: result.data[2],
                    maxMin: result.data[3],
                }
                if (this.sharesCount < this.config.minShares) this.sharesCount = this.config.minShares;
                if (this.sharesCount > this.config.maxShares) this.sharesCount = this.config.maxShares;
                if (this.threshold < this.config.minMin) this.threshold = this.config.minMin;
                if (this.threshold > this.sharesCount - 1) this.threshold = this.sharesCount - 1;
            } else {
                alert(`Error getting configuration: ${result.error}`);
            }
        },
        hideMessage() {
            const rawResult = GoSSS.hide(this.message, this.sharesCount, this.threshold);
            const result = JSON.parse(rawResult);
            if (!result.error) {
                this.hide_result = result.data.join("\n");
            } else {
                alert(`Error hiding message: ${result.error}`);
            }
        },
        recoverMessage() {
            const shares = JSON.stringify(this.shares.split("\n"));
            const rawResult = GoSSS.recover(shares);
            const result = JSON.parse(rawResult);
            if (!result.error) {
                this.recovered_message = window.atob(result.data);
            } else {
                alert(`Error recovering message: ${result.error}`);
            }
        }
    },
    watch: {
        message() {
            this.getConfig();
        }
    },
    components: {
        
    }
});

app.mount('#app');
