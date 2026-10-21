<?php
class MyInterval extends DateInterval { public function __construct($s){parent::__construct($s);} }
$i=new MyInterval('P1Y2M1W3DT4H5M6S');
echo json_encode([$i->y,$i->m,$i->d,$i->h,$i->i,$i->s,$i->f,$i->invert,$i->days]),"\n";
$d=new DateTime('2024-01-01 00:00:00');
echo $d->add($i)->format('Y-m-d H:i:s'),"\n";
echo $d->sub($i)->format('Y-m-d H:i:s'),"\n";
$i->invert=1;
$d=new DateTimeImmutable('2024-01-01 00:00:00');
$changed=$d->add(new DateInterval('PT120M'));
echo $d->format('Y-m-d H:i:s'),':',$changed->format('Y-m-d H:i:s'),"\n";
echo (new DateTime('2024-01-01 00:00:00'))->add(new DateInterval('PT0S'))->format('Y-m-d H:i:s'),"\n";
try {new DateInterval('bad');}catch(Throwable $e){echo get_class($e),':',$e->getMessage(),"\n";}
