<?php
class MissingMethod {}
$o = new MissingMethod;
$o->missing($undefined);
