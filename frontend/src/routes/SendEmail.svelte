<script>
  import { emailAPI, templateAPI } from '../lib/api';
  import { onMount } from 'svelte';

  let from = '';
  let to = '';
  let subject = '';
  let htmlBody = '';
  let textBody = '';
  let useTemplate = false;
  let templates = [];
  let selectedTemplate = '';
  let templateData = {};
  let loading = false;
  let success = '';
  let error = '';

  onMount(async () => {
    try {
      const response = await templateAPI.list();
      templates = response.data;
    } catch (err) {
      console.error('Failed to load templates:', err);
    }
  });

  async function handleSubmit() {
    error = '';
    success = '';
    loading = true;

    try {
      const data = {
        from,
        to,
      };

      if (useTemplate && selectedTemplate) {
        data.template_id = selectedTemplate;
        data.template_data = templateData;
      } else {
        data.subject = subject;
        data.html_body = htmlBody;
        data.text_body = textBody;
      }

      const response = await emailAPI.send(data);
      if (response.data.success) {
        success = `Email sent successfully! ${response.data.remaining} emails remaining.`;
        // Reset form
        to = '';
        subject = '';
        htmlBody = '';
        textBody = '';
        templateData = {};
      } else {
        error = response.data.message;
      }
    } catch (err) {
      error = err.response?.data?.error || 'Failed to send email';
    } finally {
      loading = false;
    }
  }
</script>

<div class="send-email">
  <h1>Send Email</h1>

  {#if success}
    <div class="alert alert-success">{success}</div>
  {/if}

  {#if error}
    <div class="alert alert-error">{error}</div>
  {/if}

  <form on:submit|preventDefault={handleSubmit}>
    <div class="form-group">
      <label>
        <input type="checkbox" bind:checked={useTemplate} />
        Use Template
      </label>
    </div>

    <div class="form-group">
      <label for="from">From Email</label>
      <input
        id="from"
        type="email"
        bind:value={from}
        required
        placeholder="sender@yourdomain.com"
      />
      <small>Must be from a verified domain</small>
    </div>

    <div class="form-group">
      <label for="to">To Email</label>
      <input
        id="to"
        type="email"
        bind:value={to}
        required
        placeholder="recipient@example.com"
      />
    </div>

    {#if useTemplate}
      <div class="form-group">
        <label for="template">Select Template</label>
        <select id="template" bind:value={selectedTemplate} required>
          <option value="">Choose a template</option>
          {#each templates as template}
            <option value={template.id}>{template.name}</option>
          {/each}
        </select>
      </div>

      <div class="form-group">
        <label for="templateData">Template Variables (JSON)</label>
        <textarea
          id="templateData"
          bind:value={templateData}
          rows="4"
          placeholder='{"name": "John", "url": "https://example.com"}'
        ></textarea>
      </div>
    {:else}
      <div class="form-group">
        <label for="subject">Subject</label>
        <input
          id="subject"
          type="text"
          bind:value={subject}
          required
          placeholder="Email subject"
        />
      </div>

      <div class="form-group">
        <label for="htmlBody">HTML Body</label>
        <textarea
          id="htmlBody"
          bind:value={htmlBody}
          rows="10"
          placeholder="<h1>Hello!</h1><p>Your email content here</p>"
        ></textarea>
      </div>

      <div class="form-group">
        <label for="textBody">Text Body (Optional)</label>
        <textarea
          id="textBody"
          bind:value={textBody}
          rows="5"
          placeholder="Plain text version of your email"
        ></textarea>
      </div>
    {/if}

    <button type="submit" class="btn btn-primary" disabled={loading}>
      {loading ? 'Sending...' : 'Send Email'}
    </button>
  </form>
</div>

<style>
  .send-email {
    padding: 2rem 0;
  }

  h1 {
    color: #2c3e50;
    margin-bottom: 2rem;
  }

  form {
    background: white;
    padding: 2rem;
    border-radius: 8px;
    box-shadow: 0 2px 8px rgba(0,0,0,0.1);
    max-width: 800px;
  }

  .form-group {
    margin-bottom: 1.5rem;
  }

  label {
    display: block;
    margin-bottom: 0.5rem;
    color: #2c3e50;
    font-weight: 500;
  }

  input[type="text"],
  input[type="email"],
  select,
  textarea {
    width: 100%;
    padding: 0.75rem;
    border: 1px solid #ddd;
    border-radius: 4px;
    font-size: 1rem;
    box-sizing: border-box;
    font-family: inherit;
  }

  input:focus,
  select:focus,
  textarea:focus {
    outline: none;
    border-color: #667eea;
  }

  small {
    color: #7f8c8d;
    font-size: 0.85rem;
  }

  .btn {
    padding: 0.75rem 2rem;
    border: none;
    border-radius: 4px;
    font-size: 1rem;
    cursor: pointer;
    transition: background 0.3s;
  }

  .btn-primary {
    background: #667eea;
    color: white;
  }

  .btn-primary:hover:not(:disabled) {
    background: #5568d3;
  }

  .btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .alert {
    padding: 1rem;
    border-radius: 4px;
    margin-bottom: 1.5rem;
  }

  .alert-success {
    background: #d4edda;
    color: #155724;
    border: 1px solid #c3e6cb;
  }

  .alert-error {
    background: #fee;
    color: #c00;
    border: 1px solid #fcc;
  }

  input[type="checkbox"] {
    width: auto;
    margin-right: 0.5rem;
  }
</style>
