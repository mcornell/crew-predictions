<template>
  <div class="page">
    <h1 class="page-title">Rules</h1>
    <p class="rules-intro">Three ways to be wrong. Pick your poison — actually, you don't have a choice, you're in all of the contests.</p>

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
      <p class="rules-desc">Gentler scoring for the optimists. Points stack — nail all three and you're golden.</p>
      <div class="rules-table">
        <div class="rules-row rules-row--positive">
          <span class="rules-points">+1</span>
          <span class="rules-condition">Correct result — win, draw, or loss</span>
        </div>
        <div class="rules-row rules-row--positive">
          <span class="rules-points">+1</span>
          <span class="rules-condition">Exact Columbus Crew goal tally — you predicted the right number of Crew goals</span>
        </div>
        <div class="rules-row rules-row--positive">
          <span class="rules-points">+1</span>
          <span class="rules-condition">Exact opponent goal tally — you predicted the right number of goals for the other lot</span>
        </div>
      </div>
      <p class="rules-note">Points stack. Nail all three and you're at +3.</p>
    </section>

    <section class="rules-section">
      <h2 class="section-title">Grouchy™</h2>
      <p class="rules-desc">The prediction game from the most awful forum in the most awful place to talk soccer since 2000.</p>
      <div class="rules-table">
        <div class="rules-row rules-row--positive">
          <span class="rules-points">+1</span>
          <span class="rules-condition">Win by 2+ — you predicted a comfortable Crew victory and got one</span>
        </div>
        <div class="rules-row rules-row--positive">
          <span class="rules-points">+1</span>
          <span class="rules-condition">Win by 1 — a nervy Crew win and you called it</span>
        </div>
        <div class="rules-row rules-row--positive">
          <span class="rules-points">+1</span>
          <span class="rules-condition">Draw — 90 minutes of nothing and you predicted it</span>
        </div>
        <div class="rules-row rules-row--positive">
          <span class="rules-points">+1</span>
          <span class="rules-condition">Lose by 1 — a heartbreaker and you saw it coming</span>
        </div>
        <div class="rules-row rules-row--positive">
          <span class="rules-points">+1</span>
          <span class="rules-condition">Lose by 2+ — a disaster and somehow you predicted that too</span>
        </div>
        <div class="rules-row rules-row--zero">
          <span class="rules-points">0</span>
          <span class="rules-condition">Wrong category. Touch grass.</span>
        </div>
      </div>
      <p class="rules-note">Categories are based on Columbus Crew margin. Predict the right bucket and you get the point — exact score doesn't matter.</p>
    </section>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'

onMounted(() => { document.title = 'Rules — Crew Predictions' })
</script>
