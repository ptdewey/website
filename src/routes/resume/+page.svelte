<script>
  import { education, experience, skills } from "./resume.js";
</script>

<div class="resume-container">
  <h1>Resume</h1>

  <div class="section">
    <h2>Education</h2>
    {#each education as edu}
      <div class="entry">
        <div class="entry-title">{edu.institution}</div>
        <div class="entry-subtitle">{edu.location}</div>
        <div>{edu.degree}</div>
        {#if edu.additional}
          <div>{edu.additional}</div>
        {/if}
        <div class="entry-details">{edu.dates}</div>
      </div>
    {/each}
  </div>

  <div class="section">
    <h2>Job Experience</h2>
    {#each experience as job}
      <div class="entry">
        <div class="entry-title">{job.title}</div>
        <div class="entry-subtitle">
          <span class="entry-company">{job.company}</span>, {job.location}
        </div>
        {#if job.details}
          <div>{job.details}</div>
        {/if}
        <div class="entry-details">{job.dates}</div>
      </div>
    {/each}
  </div>

  <div class="section">
    <h2>Technical Skills</h2>
    <div class="skills-list">
      <div class="skills-item">
        <strong class="skills-header">Languages:</strong>
        {skills.languages}
      </div>
      <div class="skills-item">
        <strong class="skills-header">Machine Learning & AI:</strong>
        {skills.machineLearningAI}
      </div>
      <div class="skills-item">
        <strong class="skills-header">Parallel Computing:</strong>
        {skills.parallelComputing}
      </div>
      <!-- <div class="skills-item"> -->
      <!--   <strong class="skills-header">Libraries:</strong> -->
      <!--   {skills.libraries} -->
      <!-- </div> -->
      <div class="skills-item">
        <strong class="skills-header">Web Frameworks:</strong>
        {skills.frameworks}
      </div>
      <div class="skills-item">
        <strong class="skills-header">Databases:</strong>
        {skills.databases}
      </div>
      <div class="skills-item">
        <strong class="skills-header">Messaging Tools:</strong>
        {skills.messagingTools}
      </div>
      <div class="skills-item">
        <strong class="skills-header">Operating Systems:</strong>
        {skills.devTools.operatingSystems}
      </div>
      <div class="skills-item">
        <strong class="skills-header">Version Control:</strong>
        {skills.devTools.versionControl}
      </div>
      <div class="skills-item">
        <strong class="skills-header">Containerization:</strong>
        {skills.devTools.containerization}
      </div>
      <div class="skills-item">
        <strong class="skills-header">Configuration Management:</strong>
        {skills.devTools.configurationManagement}
      </div>
      <div class="skills-item">
        <strong class="skills-header">Editors:</strong>
        {skills.devTools.editors}
      </div>
      <div class="skills-item">
        <strong class="skills-header">Cloud Services:</strong>
        {skills.devTools.cloudServices}
      </div>
      <div class="skills-item">
        <strong class="skills-header">Architecture Tools:</strong>
        {skills.architectureTools}
      </div>
    </div>
  </div>

  <div class="section">
    <h2>Projects</h2>
    <p>See <a href="/projects">my projects</a></p>
  </div>
</div>

<style>
  .resume-container {
    max-width: 1000px;
    margin: auto;
  }

  .section {
    margin-bottom: 2rem;
  }

  .entry {
    margin-bottom: 1rem;
    background-color: var(--background-color-alt);
    padding: 1rem;
    border-radius: 8px;
    transition: transform 0.2s ease;
  }

  /* .entry:hover { */
  /*   transform: scale(var(--hover-scale, 1.02)); */
  /* } */

  .entry-title {
    color: var(--yellow);
    font-size: 1.2em;
    font-weight: bold;
  }

  .entry-company {
    font-weight: bold;
    color: var(--orange);
  }

  .entry-details {
    font-size: 0.9em;
  }

  .skills-list {
    display: flex;
    flex-wrap: wrap;
    gap: 1rem;
  }

  .skills-header {
    color: var(--yellow);
    font-weight: bold;
  }

  .skills-item {
    flex: 1 1 45%;
    background-color: var(--background-color-alt);
    padding: 1rem;
    border-radius: 8px;
    transition: transform 0.2s ease;
  }

  /* .skills-item:hover { */
  /*   transform: scale(var(--hover-scale, 1.02)); */
  /* } */

  @media (max-width: 768px) {
    .skills-item {
      flex: 1 1 100%;
    }
  }
</style>
