<template>
  <div class="page">
    <h1 class="page-title">Rules</h1>
    <p class="rules-intro">Two ways to be wrong. Pick your poison.</p>

    <section class="rules-section">
      <h2 class="section-title">Aces Radio</h2>
      <p class="rules-desc">High risk, high reward. Nail the exact scoreline and you're the smartest person in the room. Get it backwards and you owe everyone a pint.</p>
      <div class="rules-table">
        <div class="rules-row rules-row--positive">
          <span class="rules-points">+15</span>
          <span class="rules-condition">Exact scoreline — you called it perfectly</span>
        </div>
        <div class="rules-row rules-row--positive">
          <span class="rules-points">+10</span>
          <span class="rules-condition">Correct winner or draw, wrong score</span>
        </div>
        <div class="rules-row rules-row--negative">
          <span class="rules-points">-15</span>
          <span class="rules-condition">Flipped — you predicted the exact same score but reversed (e.g. you said 2-1 Crew win, actual was 1-2 loss)</span>
        </div>
        <div class="rules-row rules-row--zero">
          <span class="rules-points">0</span>
          <span class="rules-condition">Everything else. Wrong result, wrong score, wrong life choices.</span>
        </div>
      </div>
    </section>

    <section class="rules-section">
      <h2 class="section-title">Upper 90 Club</h2>
      <p class="rules-desc">Gentler scoring for the optimists. Points stack — get both right and you're golden.</p>
      <div class="rules-table">
        <div class="rules-row rules-row--positive">
          <span class="rules-points">+1</span>
          <span class="rules-condition">Correct result — win, draw, or loss</span>
        </div>
        <div class="rules-row rules-row--positive">
          <span class="rules-points">+1</span>
          <span class="rules-condition">Exact Columbus Crew goal tally — you predicted the right number of Crew goals</span>
        </div>
      </div>
      <p class="rules-note">Points stack. Nail the result and the Crew's goal tally and you're at +2.</p>
    </section>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'

onMounted(() => { document.title = 'Rules — Crew Predictions' })
</script>
