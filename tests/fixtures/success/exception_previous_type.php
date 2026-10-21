<?php
try { new Exception(previous: new stdClass); }
catch (TypeError $error) { echo $error->getMessage(), "\n"; }
