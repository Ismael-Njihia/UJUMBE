<script>
  import { onMount } from 'svelte';
  import { templateAPI } from '../lib/api';

  let templates = [];
  let loading = true;
  let showCreateForm = false;
  let name = '';
  let subject = '';
  let htmlBody = '';
  let textBody = '';
  let error = '';
  let success = '';

  onMount(loadTemplates);

  async function loadTemplates() {
    loading = true;
    try {
      const response = await templateAPI.list();
      templates = response.data || [];
    } catch (err) {
      error = 'Failed to load templates';
    } finally {
      loading = false;
    }
  }

  async function handleCreate() {
    error = '';
    success = '';

    try {
      await templateAPI.create({ name, subject, html_body: htmlBody, text_body: textBody });
      success = 'Template created successfully!';
      showCreateForm = false;
      name = '';
      subject = '';
      htmlBody = '';
      textBody = '';
      await loadTemplates();
    } catch (err) {
      error = err.response?.data?.error || 'Failed to create template';
    }
  }

  async function handleDelete(id) {
    if (!confirm('Are you sure you want to delete this template?')) return;

    try {
      await templateAPI.delete(id);
      success = 'Template deleted successfully!';
      await loadTemplates();
    } catch (err) {
      error = 'Failed to delete template';
    }
  }
</script>

<div class="templates">
  <div class="header">
    <h1>Email Templates</h1>
    <button class="btn btn-primary" on:click={() => showCreateForm = !showCreateForm}>
      {showCreateForm ? 'Cancel' : 'Create Template'}
    </button>
  </div>

  {#if success}
    <div class="alert alert-success">{success}</div>
  {/if}

  {#if error}
    <div class="alert alert-error">{error}</div>
  {/if}

  {#if showCreateForm}
    <div class="create-form">
      <h2>Create New Template</h2>
      <form on:submit|preventDefault={handleCreate}>
        <div class="form-group">
          <label for="name">Template Name</label>
          <input id="name" type="text" bind:value={name} required placeholder="Welcome Email" />
        </div>

        <div class="form-group">
          <label for="subject">Subject</label>
          <input id="subject" type="text" bind:value={subject} required placeholder="Welcome to {{name}}" />
        </div>

        <div class="form-group">
          <label for="htmlBody">HTML Body</label>
          <textarea
            id="htmlBody"
            bind:value={htmlBody}
            required
            rows="10"
            placeholder="<h1>Welcome {{name}}!</h1>"
          ></textarea>
          <small>Use {"{{variable}}"} for dynamic content</small>
        </div>

        <div class="form-group">
          <label for="textBody">Text Body (Optional)</label>
          <textarea
            id="textBody"
            bind:value={textBody}
            rows="5"
            placeholder="Welcome {{name}}!"
          ></textarea>
        </div>

        <button type="submit" class="btn btn-primary">Create Template</button>
      </form>
    </div>
  {/if}

  {#if loading}
    <p>Loading templates...</p>
  {:else if templates.length === 0}
    <div class="empty-state">
      <p>No templates yet. Create your first template to get started!</p>
    </div>
  {:else}
    <div class="templates-grid">
      {#each templates as template}
        <div class="template-card">
          <h3>{template.name}</h3>
          <p class="subject"><strong>Subject:</strong> {template.subject}</p>
          <div class="preview">
            {@html template.html_body.substring(0, 100) + (template.html_body.length > 100 ? '...' : '')}
          </div>
          <div class="actions">
            <button class="btn btn-small btn-danger" on:click={() => handleDelete(template.id)}>
              Delete
            </button>
          </div>
        </div>
      {/each}
    </div>
  {/if}
</div>

<style>
  .templates {
    padding: 2rem 0;
  }

  .header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 2rem;
  }

  h1 {
    color: #2c3e50;
    margin: 0;
  }

  .create-form {
    background: white;
    padding: 2rem;
    border-radius: 8px;
    box-shadow: 0 2px 8px rgba(0,0,0,0.1);
    margin-bottom: 2rem;
  }

  .create-form h2 {
    margin-top: 0;
    color: #2c3e50;
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

  input,
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
  textarea:focus {
    outline: none;
    border-color: #667eea;
  }

  small {
    color: #7f8c8d;
    font-size: 0.85rem;
  }

  .templates-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: 1.5rem;
  }

  .template-card {
    background: white;
    padding: 1.5rem;
    border-radius: 8px;
    box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  }

  .template-card h3 {
    margin: 0 0 1rem 0;
    color: #2c3e50;
  }

  .subject {
    color: #7f8c8d;
    font-size: 0.9rem;
    margin-bottom: 1rem;
  }

  .preview {
    color: #95a5a6;
    font-size: 0.85rem;
    margin-bottom: 1rem;
    padding: 1rem;
    background: #f8f9fa;
    border-radius: 4px;
    max-height: 100px;
    overflow: hidden;
  }

  .actions {
    display: flex;
    gap: 0.5rem;
  }

  .btn {
    padding: 0.75rem 1.5rem;
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

  .btn-primary:hover {
    background: #5568d3;
  }

  .btn-small {
    padding: 0.5rem 1rem;
    font-size: 0.9rem;
  }

  .btn-danger {
    background: #e74c3c;
    color: white;
  }

  .btn-danger:hover {
    background: #c0392b;
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

  .empty-state {
    background: white;
    padding: 3rem;
    border-radius: 8px;
    text-align: center;
    color: #7f8c8d;
  }
</style>
