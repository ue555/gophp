<?php
function emptyBody() {}
function emptyReturn() { return; echo 'bad'; }
echo (emptyBody() === null) + 0, (emptyReturn() === null) + 0, "\n";
