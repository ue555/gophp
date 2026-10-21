<?php
echo hash('xxh128', ''), '|', hash('xxh128', 'abc'), '|', strlen(hash('xxh128', 'abc', true));
