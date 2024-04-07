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
        <div style="width: 90%; max-width: 800px; margin: 50px auto;" class="is-shadowed is-rounded has-p-12">
            <header>
                <h1 class="title">Shamir's Secret Sharing Demo</h1>
                <p>Welcome to the Shamir's Secret Sharing demo 👋.</p> 
                <p>This tool allows you to securely <b>share a secret message by dividing it</b> into parts. A <b>certain number</b> of parts (threshold) <b>are needed to recover</b> the original message. </p>
                <p>This ensures that the secret can only be reconstructed when a sufficient number of parts are brought together.</p>
            </header>

            <div class="is-flex has-direction-row has-text-center">
                <button class="button has-m-2 has-w-full" :class="{'is-normal': currentTab != 'hide'}" @click="currentTab = 'hide'">Hide</button>
                <button class="button has-m-2 has-w-full" :class="{'is-normal': currentTab != 'recover'}" @click="currentTab = 'recover'">Recover</button>
            </div>
            <div v-show="currentTab === 'hide'">
                <h3 class="has-mt-6 has-mb-8">Hide a message</h3>
                <textarea class="textarea has-mb-4" v-model="message" placeholder="Enter your message" rows="6"></textarea>
                <div class="is-flex has-direction-row has-justify-center has-mb-4">
                    <div class="has-w-full has-m-2" v-show="message">
                        <label class="label">Shares </label>
                        <input v-model="sharesCount" type="number" class="input" :min="config.minShares" :max="config.maxShares" :step="config.minShares">
                        <small>(min {{ config.minShares }}, max {{ config.maxShares }})</small>
                    </div>
                    <div class="has-w-full has-m-2" v-show="message">
                        <label class="label">Threshold</label>
                        <input v-model="threshold" type="number" class="input" :min="config.minMin" :max="config.maxMin">
                        <small>(min {{ config.minMin }}, max {{ config.maxMin }})</small>
                    </div>
                </div>
                <button class="button is-secondary has-w-full has-mt-4 has-mb-4" :class="{'is-disabled': !message}" @click="hideMessage">Hide</button>
                <div class="has-mt-4 has-mb-4" v-show="hide_result">
                    <h4 class="has-mt-4 has-mb-6">Resulting secret parts</h4>
                    <textarea class="textarea" readonly v-model="hide_result" placeholder="Resulting secret parts" rows="6"></textarea>
                </div>
            </div>
            <div v-show="currentTab === 'recover'">
                <h3 class="has-mt-6 has-mb-8">Recover a message</h3>
                <textarea class="textarea has-mb-4" v-model="shares" placeholder="Enter shares" rows="6"></textarea>
                <button class="button is-secondary has-w-full has-mt-4 has-mb-4" :class="{'is-disabled': !shares}" @click="recoverMessage">Recover</button>
                <div class="has-mt-4 has-mb-4" v-show="recovered_message">
                    <h4 class="has-mt-4 has-mb-6">Recovered message</h4>
                    <textarea class="textarea" readonly v-model="recovered_message" placeholder="Recovered message" rows="6"></textarea>
                </div>
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
                    minShares: result.data.minShares,
                    maxShares: result.data.maxShares,
                    minMin: result.data.minMin,
                    maxMin: result.data.maxMin,
                }
                if (this.sharesCount < this.config.minShares) this.sharesCount = this.config.minShares;
                if (this.sharesCount > this.config.maxShares) this.sharesCount = this.config.maxShares;
                this.config.maxMin = this.sharesCount - (this.config.minShares / 3);
                if (this.threshold < this.config.minMin) this.threshold = this.config.minMin;
                if (this.threshold > this.config.maxMin) this.threshold = this.config.maxMin;
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
        },
        sharesCount() {
            this.getConfig();
        },
    },
});

app.mount('#app');
