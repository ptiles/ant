# Running Langton's ant on Penrose tiling

## Building

```shell
make
```

## Running

For full usage, run:

```shell
./bin/ant-batch -h
```

```shell
./bin/ant -h
```

## Batch run

```shell
./bin/ant-batch -n LLR,LRR,RLL,RRL -ic 4 -rc 4 -d results_4_4_500k -s 500_000
```

![results_4_4_500k.png](example/results_4_4_500k.png)

## Example results

```shell
./bin/ant LRR__0.641800__D2139+C-8152__2_500_000
```

![LRR__0.641800__D2139+C-8152__2500000.png](example/LRR__0.641800__D2139%2BC-8152__2500000.png)

```shell
./bin/ant LRR__0.246000__A7386-B5868__3_500_000_000
```

![LRR__0.246000__A7386-B5868__3500000000.png](example/LRR__0.246000__A7386-B5868__3500000000.png)

## References

- [Langton's Ant on Penrose Tiling producing a Pentaflake-like Fractal (video) by dropped box](https://www.youtube.com/watch?v=vUdfcftF5cM)
- [Langton's Ant draws Koch Snowflake on Penrose Tiling (video) by dropped box](https://www.youtube.com/watch?v=D72Op1Z_VFQ)
- [Pattern Collider by Aatish Bhatia, Henry Reich](https://aatishb.com/patterncollider/)
- [Pattern Collider (GitHub)](https://github.com/aatishb/patterncollider)
- [Pentagrids and Penrose Tilings by Stacy Mowry, Shriya Shukla](https://web.williams.edu/Mathematics/sjmiller/public_html/hudson/HRUMC-Mowry&Shukla_Pentagrids%20and%20Penrose.pdf)
- [deBruijn Mathematical Details by Greg Egan](https://www.gregegan.net/APPLETS/12/deBruijnNotes.html)
