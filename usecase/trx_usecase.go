package usecase

import (
	"errors"
	"fmt"
	"gogroceries/domain"          
	"gogroceries/internal/helper" 
	"log"
	"math"

	"gorm.io/gorm"
)

type trxUsecase struct {
	trxRepo      domain.TrxRepository
	produkRepo   domain.ProdukRepository
	alamatRepo  domain.AlamatRepository 
	categoryRepo domain.CategoryRepository 
	tokoRepo    domain.TokoRepository    
}

func NewTrxUsecase(
    tr domain.TrxRepository,
    pr domain.ProdukRepository,
    ar domain.AlamatRepository,
    cr domain.CategoryRepository,
    trRepo domain.TokoRepository,
) domain.TrxUsecase {
    return &trxUsecase{
        trxRepo:      tr,
        produkRepo:   pr,
        alamatRepo:   ar,
        categoryRepo: cr,
        tokoRepo:     trRepo,
    }
}

func (uc *trxUsecase) CreateTransaksi(req *domain.CreateTransaksiRequest, userID uint) (*domain.Trx, error) {
	alamat, err := uc.alamatRepo.FindByIDAndUserID(req.IdAlamatKirim, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("alamat pengiriman tidak valid atau bukan milik anda")
		}
		return nil, fmt.Errorf("gagal validasi alamat: %w", err)
	}

	var detailsToSave []domain.DetailTrx
	var logsToSave []domain.LogProduk
	var totalHargaKeseluruhan int = 0

	productIDs := make([]uint, len(req.DetailTrx))
	for i, item := range req.DetailTrx {
		productIDs[i] = item.IdProduk
	}

	produks, err := uc.produkRepo.FindByIDs(productIDs) 
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data produk: %w", err)
	}

	produkMap := make(map[uint]*domain.Produk)
	for _, p := range produks {
		produkMap[p.ID] = &p
	}

	for _, item := range req.DetailTrx {
		produk, ok := produkMap[item.IdProduk]
		if !ok {
			return nil, fmt.Errorf("produk dengan ID %d tidak ditemukan", item.IdProduk)
		}

		if produk.Stok < item.Kuantitas {
			return nil, fmt.Errorf("stok produk '%s' tidak mencukupi (tersedia: %d, diminta: %d)", produk.NamaProduk, produk.Stok, item.Kuantitas)
		}

		if produk.Toko == nil || produk.Category == nil {
			log.Printf("Fatal Error: Gagal preload Toko/Category untuk produk ID %d", item.IdProduk)
			return nil, fmt.Errorf("gagal mendapatkan detail toko/kategori untuk produk ID %d", item.IdProduk)
		}

		hargaItemTotal := produk.HargaKonsumen * item.Kuantitas
		totalHargaKeseluruhan += hargaItemTotal

		detailsToSave = append(detailsToSave, domain.DetailTrx{
			IdProduk:   produk.ID,
			IdToko:     produk.IdToko,
			Kuantitas:  item.Kuantitas,
			HargaTotal: hargaItemTotal,
		})

		logsToSave = append(logsToSave, domain.LogProduk{
			NamaProduk:    produk.NamaProduk,
			Slug:          produk.Slug,
			HargaReseller: produk.HargaReseller,
			HargaKonsumen: produk.HargaKonsumen,
			Deskripsi:     produk.Deskripsi,
			IdCategory:    produk.IdCategory,
			NamaCategory:  produk.Category.NamaCategory,
			IdToko:        produk.IdToko,
			NamaToko:      produk.Toko.NamaToko,
			UrlFotoToko:   produk.Toko.UrlFoto,
		})
		
	}

	kodeInvoice := helper.GenerateInvoiceCode()

	newTrx := &domain.Trx{
		IdUser:        userID,
		IdAlamatKirim: alamat.ID,
		HargaTotal:    totalHargaKeseluruhan,
		KodeInvoice:   kodeInvoice,
		MethodBayar:   req.MethodBayar,
		Status:        "pending",
	}

	createdTrx, err := uc.trxRepo.Create(newTrx, detailsToSave, logsToSave)
	if err != nil {
		return nil, err
	}

	return createdTrx, nil
}


func (uc *trxUsecase) GetAllTransaksiUser(userID uint, filter domain.TrxFilter, page, limit int) ([]domain.Trx, *domain.PaginationResponse, error) {
    if page < 1 {
        page = 1
    }
    if limit < 1 {
        limit = 10 
    }
    offset := (page - 1) * limit

    trxs, totalData, err := uc.trxRepo.FindAllByUserID(userID, filter, limit, offset)
    if err != nil {
        return nil, nil, err
    }

    totalPage := int(math.Ceil(float64(totalData) / float64(limit)))
    
    pagination := &domain.PaginationResponse{
        Page:      page,
        Limit:     limit,
        TotalData: int(totalData),
        TotalPage: totalPage,
        Data:      trxs,
    }

    return trxs, pagination, nil
}

func (uc *trxUsecase) GetTransaksiByID(id uint, userID uint) (*domain.Trx, error) {
	trx, err := uc.trxRepo.FindByID(id, userID) 
	if err != nil {
		return nil, err
	}
	return trx, nil
}