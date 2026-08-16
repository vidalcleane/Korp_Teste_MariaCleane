import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { MatTableModule } from '@angular/material/table';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { NotaService, ItemNota } from '../services/nota';
import { ProdutoService, Produto } from '../services/produto';

@Component({
  selector: 'app-notas',
  imports: [
    CommonModule,
    ReactiveFormsModule,
    MatTableModule,
    MatFormFieldModule,
    MatInputModule,
    MatSelectModule,
    MatButtonModule,
    MatCardModule,
    MatIconModule,
    MatProgressSpinnerModule,
    MatSnackBarModule
  ],
  templateUrl: './notas.html',
  styleUrl: './notas.css',
})
export class Notas implements OnInit {
  colunas = ['numero', 'status', 'itens', 'acoes'];
  itemForm: FormGroup;
  carrinho: ItemNota[] = [];

  // Controla quais notas estão sendo impressas agora, para mostrar o
  // indicador de processamento (spinner) só no botão certo.
  imprimindo = new Set<number>();

  produtos: Produto[] = [];

  constructor(
    private notaService: NotaService,
    private produtoService: ProdutoService,
    private fb: FormBuilder,
    private snackBar: MatSnackBar
  ) {
    this.itemForm = this.fb.group({
      produto_id: ['', Validators.required],
      quantidade: [1, [Validators.required, Validators.min(1)]],
    });
  }

  ngOnInit(): void {
    this.notaService.carregarNotas();
    this.produtoService.carregarProdutos();
    this.produtoService.produtos$.subscribe((lista) => (this.produtos = lista));
  }

  get notas$() {
    return this.notaService.notas$;
  }

  nomeProduto(id: number): string {
    const produto = this.produtos.find((p) => p.id === id);
    return produto ? produto.codigo : `#${id}`;
  }

  adicionarItem(): void {
    if (this.itemForm.invalid) return;
    this.carrinho.push(this.itemForm.value);
    this.itemForm.reset({ produto_id: '', quantidade: 1 });
  }

  removerItem(index: number): void {
    this.carrinho.splice(index, 1);
  }

  criarNota(): void {
    if (this.carrinho.length === 0) {
      this.snackBar.open('Adicione ao menos um produto à nota.', 'Fechar', { duration: 3000 });
      return;
    }

    this.notaService.criarNota(this.carrinho).subscribe({
      next: () => {
        this.snackBar.open('Nota fiscal criada com sucesso!', 'Fechar', { duration: 3000 });
        this.carrinho = [];
      },
      error: () => {
        this.snackBar.open('Erro ao criar a nota fiscal.', 'Fechar', { duration: 4000 });
      }
    });
  }

  imprimir(numero: number): void {
    this.imprimindo.add(numero);

    this.notaService.imprimirNota(numero).subscribe({
      next: () => {
        this.imprimindo.delete(numero);
        this.snackBar.open(`Nota ${numero} impressa com sucesso!`, 'Fechar', { duration: 3000 });
      },
      error: (err) => {
        this.imprimindo.delete(numero);
        // Aqui aparece o feedback do cenário de falha: se o Estoque
        // estiver fora do ar, a mensagem específica do backend chega até aqui.
        const mensagem = err.error?.erro || 'Erro ao imprimir a nota.';
        this.snackBar.open(mensagem, 'Fechar', { duration: 6000 });
      }
    });
  }
}
