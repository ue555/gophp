<?php
class C {}
$o = new C; echo ++$o->n, ':', $o->n++, ':', $o->n, "\n";
