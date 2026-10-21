<?php
echo phpversion(), '|';
echo phpversion('unknown') === false ? 'missing' : 'bad';
