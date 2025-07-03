---
sidebar_position: 4
title: Score Aggregation
description: How rankings are collected from judges and scores are aggregated.
---

# Score Aggregation

Jury's main feature revolves around its rank aggregation system. As judges view projects, they are asked to provide input through ranking and starring projects. These metrics will be aggregated to form a comprehensive ranking. In the following section, we go more in-depth on how this is done.

## Ranking

The basis of Jury's algorithm is the [Borda Count](https://en.wikipedia.org/wiki/Borda_count) method. A judge will rank up to 5 projects they have seen:

![Judge ranking](./assets/judge-ranking.png)

For each of the projects that they rank, the first place ranking will be assigned 5 points; the second place will be assigned 4 points; and so on. If, for example, a judge only ranks 3 projects, then the first place ranking will be assigned only 3 points. The points from all judges will be tallied up, resulting in the final score.

This method of ranking was chosen for its simplicity yet accurate aggregated ranking, [with similar error rates to other ranking algorithms](https://arxiv.org/pdf/2204.05742).

## Starring

Projects can be starred from the finish judging popup: 

![finish judging popup](./assets/judge-popup-star.png)

Or from the ranking menu:

![ranking menu](./assets/ranking-star.png)

Stars are simply added up for each judge. We introduced stars as a way for judges to pick their absolute favorites out of the projects that they've seen. If they saw a project that they think should win the hackathon, we encourage judges to star those projects. There is no limit to how many projects a judge can star, but we recommend they star at most 3. The star system allows for judges to provide more information to the organizers and allows for organizers to get a better sense of which projects were most loved. We recommend using the score system mainly but considering the number of stars as a secondary metric.

## Final Deliberation

Once all projects are judged, organizers will go into a room to do final deliberation. They will use the scores from Jury to determine the final ranking of projects, as well as the rankings among any tracks. To help in the final decision, it's a good idea to send organizers out before final scores are tallied and view the top ranked projects on Jury, especially if there are multiple that stand out among the best.
