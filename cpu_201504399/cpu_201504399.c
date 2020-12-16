#include <linux/proc_fs.h>
#include <linux/kernel.h>
#include <linux/module.h>
#include <linux/init.h>
#include <linux/sched/signal.h>
#include <linux/seq_file.h>
#include <linux/sched.h>
#include <linux/fs.h>
#include <linux/mm.h>


struct task_struct *task; //estructura definifa en sched.h para tareas /procesos 
struct task_struct *task_child; //estructura necesaria para iterar a través de tareas secundarias
struct list_head *list;  //estructura necesaria para recorrer la lista en cada tarea -> estructura de hijos.

static int escribir_archivo(struct seq_file * archivo, void *v){
    int contador =0;
    int contadoraux=0;
     #define Convert(x) ((x) << (PAGE_SHIFT - 10))
    seq_printf(archivo, " [ ");

    for_each_process(task){
        if(contador>0){
            seq_printf(archivo,"   , \n");
        }
        seq_printf(archivo," {\n");
        seq_printf(archivo, "\"PID\": \"%d\",",task->pid);
        seq_printf(archivo, "\n\"PROCESO\": \"%s\",", task->comm);
        seq_printf(archivo, "\n\"ESTADO\": \"%ld\",",task->state);
        if(task->mm!=NULL){
            seq_printf(archivo, "\n\"MEMORIA\": \"%ld\",",Convert(get_mm_rss(task->mm)));
        } else {
            seq_printf(archivo,"\n\"MEMORIA\":\"0\",");
        }        
        seq_printf(archivo,"\n\"USUARIO\": \"%d\",", __kuid_val(task->real_cred->uid));
    
        seq_printf(archivo, "\n\"HIJOS\": [ ");
        contadoraux=0;
        list_for_each(list, &task->children){
            if(contadoraux>0){
                seq_printf(archivo,"   , \n");
            }
            task_child = list_entry(list,struct task_struct, sibling);
            seq_printf(archivo,"\n\t{\n");
            seq_printf(archivo,"\t\"PID\": \"%d\",",task_child->pid);
            seq_printf(archivo,"\n\t\"PROCESO\": \"%s\",",task_child->comm);
            seq_printf(archivo,"\n\t\"ESTADO\": \"%ld\",",task_child->state);
            if(task_child->mm!=NULL){
                    seq_printf(archivo, "\n\t\"MEMORIA\":\"%ld\",",Convert(get_mm_rss(task_child->mm)));
            } else {
                    seq_printf(archivo,"\n\t\"MEMORIA\":\"0\",");
            }
            seq_printf(archivo,"\n\t\"USUARIO\": \"%d\"", __kuid_val(task_child->real_cred->uid));
            seq_printf(archivo,"\n\t } \n");
            contadoraux++;
        }
        seq_printf(archivo," ] \n");
        contador++;
        seq_printf(archivo,"}  \n");

    }
    seq_printf(archivo, " ] ");
    #undef k
    return 0;
}

static int al_abrir(struct inode *inode, struct file *file){
    return single_open(file,escribir_archivo,NULL); 
}

static struct file_operations operaciones = 
{
    .open = al_abrir,
    .read = seq_read
};

int iniciar(void){ //modulo de inicio 
    proc_create("cpu_201504399",0,NULL,&operaciones);
    printk(KERN_INFO "%s","Cargando modulo.\n");
    printk(KERN_INFO "%s","María de los Angeles Herrera Sumalé.\n"); 
   /* for_each_process(task){

        printk(KERN_INFO "\nPADRE PID: %d PROCESO: %s ESTADO: %ld",task->pid, task->comm,task->state);
        list_for_each(list, &task->children){
            task_child = list_entry(list,struct task_struct, sibling);
            printk(KERN_INFO "\n HIJO DE %s[%d] PID: %d PROCESO: %s  ESTADO:  %ld",task->comm, task->pid,task_child->pid,
            task_child->comm,task_child->state);
            
        }
        printk("*****************************************************");
    } 
    printk(KERN_INFO "%s","Cargando modulo.\n");
    printk(KERN_INFO "%s","María de los Angeles Herrera Sumalé.\n"); */
    return 0;
}

void salir(void){ //modulo salida 
    remove_proc_entry("cpu_201504399",NULL);
    printk(KERN_INFO "%s","Removiendo modulo.\n");
    printk(KERN_INFO "Diciembre 2020.\n");
}

module_init(iniciar);
module_exit(salir);

MODULE_LICENSE("GPL");
MODULE_DESCRIPTION("Itera a través de todos los procesos / procesos hijos en el SO.");
MODULE_AUTHOR("201504399");
