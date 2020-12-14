(load "mk-vicare.scm")
(load "mk.scm")
(import (Framework match))
(import (Framework helpers))
(case-sensitive #t)


(define MAX_NUM 10000000)

(define (alwaysᵒ)
  (conde
   [succeed]
   [(alwaysᵒ)]))

(define (caro p a)
  (fresh (d)
    (== p (cons a d))))

(define (cdro p d)
  (fresh (a)
    (== p (cons a d))))

(define (memberᵒ x l)
  (conde
   [(caro l x)]
   [(fresh (d)
      (cdro l d)
      (memberᵒ x d))]))


(define (nullᵒ l)
  (== l '()))


(define (lorᵒ relᵒ l)
  (conde
   [(nullᵒ l)]
   [(fresh (a)
      (caro l a)
      (relᵒ a))
    (fresh (d)
      (cdro l d)
      (lorᵒ relᵒ d))]))

(define (random-outᵒ r out)
  (conde
   [(== r (+ out (random (* 10 out))))]
   [(random-outᵒ r out)]))


(define (random-nᵒ n out)
  (conde
   [(== out (random n))]
   [(random-nᵒ n out)]))


(run 10 (r)
  (random-nᵒ 10000 r))


(define (randᵒ n out)
  (cond
    [(and (number? n) (number? out)) (if (< out n) succeed fail)]
    [(number? n) (random-nᵒ n out)]
    [(number? out) (random-outᵒ n out)]
    [else
     (let ([rn  (random MAX_NUM)]
           [rout (random MAX_NUM)])
       (conde
        [(== n rn) (random-nᵒ rn out)]
        [(== out rout) (random-outᵒ n rout)]))]))

(define (randsymᵒ n)
  (conde
   [(== n (gensym))]
   [(randsymᵒ n)]))

(define (randsym-may-dupᵒ n)
  (conde
   [(alwaysᵒ) (== n (gensym))]
   [(randsym-may-dupᵒ n)]))

(define (set-stmt⁰ stmt)
  (fresh (s n)
    (== stmt `(set! ,s ,n))
    (randsymᵒ s)
    (randomᵒ n)))

(define (predᵒ p)
  (conde
   [(== p #t)]
   [(== p #f)]
   [(fresh (n1 n2)
      (== p `(> ,n1 ,n2))
      (randᵒ 100000 n1)
      (randᵒ 100000 n2))]
   [(fresh (pred conseq alter)
      (== p `(if ,pred ,conseq ,alter))
      (predᵒ pred)
      (predᵒ conseq)
      (predᵒ alter))]))



(define (Registerᵒ r)
  (memberᵒ r registers))


(define (Registerᵒ r)
  (conde
   [(alwaysᵒ)
    (memberᵒ r registers)]
   [(Registerᵒ r)]))

(define (Registerᵒ r)
  (let* ([n (random (length registers))]
         [reg (list-ref registers n)])
   (conde
    [(alwaysᵒ) (== r reg)]
    [(Registerᵒ r)])))

(define Varᵒ Registerᵒ)

(define (Binopᵒ b)
  (conde
   [(alwaysᵒ)
    (memberᵒ b '(+ - *))]
   [(Binopᵒ b)]))

(define (Binopᵒ r)
  (let* ([l '(+ - *)]
         [n (random (length l))]
         [reg (list-ref l n)])
   (conde
    [(alwaysᵒ) (== r reg)]
    [(Binopᵒ r)])))


(define (Statementᵒ s)
  (conde
   [(fresh (v1 v2 i b)
      (Varᵒ v1)
      (Varᵒ v2)
      (randᵒ 10000 i)
      (Binopᵒ b)
      (conde
       [(== s `(set! ,v1 ,i))]
       [(== s `(set! ,v1 ,v2))]
       [(== s `(set! ,v1 (,b ,v1 ,i)))]
       [(== s `(set! ,v1 (,b ,v1 ,v2)))]
       [(fresh (s1 s*)         
          (== s `(begin ,s1 ,@s*))
          (Statementᵒ s1)
          (lorᵒ Statementᵒ s*))]))]))

(run 100 (q)
  (Statementᵒ q))

#!eof
(load "handy.ss")

(run 20 (q)
      (conde
       [(alwaysᵒ) (memberᵒ q '(rax rbx rcx))]))

(car
 (run 5 (q)
    (random-nᵒ 1000 q)))

(run 5 (r)
  (random-outᵒ r 1000))

(run 10 (s)
  (randsymᵒ s))

(run 10 (s)
  (set-stmt⁰ s))


(run 500 (p)
  (carᵒ p 'if)
  (predᵒ p))
