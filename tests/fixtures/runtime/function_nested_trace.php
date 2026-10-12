<?php
function inner($n) { $n += 1; return 1 / 0; }
function outer($n) { return inner($n); }
outer(3);
