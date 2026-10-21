<?php
echo strip_tags('<b>A</b> &amp; <i>B</i>'), '|';
echo html_entity_decode('A &amp; B'), '|';
echo mb_strimwidth('Hello World', 0, 10, '...'), '|';
echo mb_strimwidth('日本語', 0, 4), '|';
usleep(0);
echo 'done';
